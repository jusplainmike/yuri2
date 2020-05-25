package restapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zekroTJA/yuri2/internal/auth"
	"github.com/zekroTJA/yuri2/internal/storage"
	"github.com/zekroTJA/yuri2/internal/util"
)

var (
	jwtCookieKey = "_session_jwt"
)

func (r *RestAPI) hAuthLogin(ctx *gin.Context) {
	r.doa.HandleInitialize(ctx.Writer, ctx.Request)
	ctx.Abort()
}

func (r *RestAPI) hAuthCallback(ctx *gin.Context) {
	user, err := r.doa.HandleCallback(ctx.Writer, ctx.Request)
	if r.doa.IsErrUnauthorized(err) {
		failUnauthorized(ctx)
		return
	}
	if err != nil {
		failInternal(ctx, err)
		return
	}

	token, err := r.auth.GenerateSessionKey(user)
	if err != nil {
		failInternal(ctx, err)
		return
	}

	ctx.SetCookie(jwtCookieKey, token, int(r.auth.GetExpireTime().Seconds()), "", "", false, true)

	ok(ctx)
}

func (r *RestAPI) hAuthValidate(ctx *gin.Context) {
	rak := ctx.Query("rak")
	if rak != "" {
		if rak != r.rak {
			failUnauthorized(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
		return
	}

	token, _ := ctx.Cookie(jwtCookieKey)
	if token == "" {
		failUnauthorized(ctx)
		ctx.Abort()
		return
	}

	data, err := r.auth.ValidateSessionKey(token)
	if err == auth.ErrInvalidSessionKey {
		failUnauthorized(ctx)
		ctx.Abort()
		return
	} else if err != nil {
		failInternal(ctx, err)
		ctx.Abort()
		return
	}

	userId, ok := data.(string)
	if !ok {
		failInternal(ctx, errors.New("invalid token data type"))
		ctx.Abort()
		return
	}

	ctx.Set("userId", userId)
	ctx.Next()
}

func (r *RestAPI) hSoundsPost(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		failBadRequest(ctx)
		return
	}

	fh, err := file.Open()
	if err != nil {
		failInternal(ctx, err)
		return
	}
	defer fh.Close()

	mimeType := util.GetMimeType(file.Filename)

	fileInfo := &storage.File{
		ContentType: mimeType,
		Modified:    time.Now(),
		Name:        file.Filename,
		Size:        file.Size,
	}

	if err = checkFile(fileInfo); err != nil {
		fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = r.s.Put(file.Filename, fh, file.Size, mimeType)
	if err != nil {
		failInternal(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, fileInfo)
}

func (r *RestAPI) hSoundsGet(ctx *gin.Context) {
	limit := util.AtoiDef(ctx.Query("limit"), 1000)

	files := make([]*storage.File, limit)

	cdone := make(chan struct{})
	cfiles := r.s.ListAsync(cdone)

	var i int
	for f := range cfiles {
		if f.Err != nil {
			failInternal(ctx, f.Err)
			return
		}

		files[i] = f

		i++
		if i >= limit {
			cdone <- struct{}{}
			break
		}
	}

	files = files[0:i]
	ctx.JSON(http.StatusOK, &listResponse{
		Size: i,
		Data: files,
	})
}

func (r *RestAPI) hSoundGet(ctx *gin.Context) {
	fileName := ctx.Param("fileName")

	mimeType := util.GetMimeType(fileName)

	file, err := r.s.Get(fileName)
	if r.s.IsNotFoundErr(err) {
		failNotFound(ctx)
		return
	} else if err != nil {
		failInternal(ctx, err)
		return
	}
	defer file.Close()

	ctx.DataFromReader(http.StatusOK, file.Size, mimeType, file.Reader, nil)
}
