package restapi

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/zekroTJA/yuri2/internal/storage"
	"github.com/zekroTJA/yuri2/internal/util"
)

var (
	allowedMimeTypes       = []string{"audio/ogg", "audio/mpeg"}
	allowedFileName        = regexp.MustCompile(`^[a-z0-9]+\.[a-z0-9]{3,5}$`)
	allowedFileSize  int64 = 100 << 20 // 100 MB

	errInvalidMimeType = fmt.Errorf("unallowed mime type - must be one of %v", allowedMimeTypes)
	errInvalidFileName = fmt.Errorf("invalid file name - file name must match %s", allowedFileName.String())
	errInvalidFileSize = fmt.Errorf("invalid file size - must not be larger than %d", allowedFileSize)
)

func fail(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"code":  code,
		"error": message,
	})
}

func failNotFound(ctx *gin.Context) {
	fail(ctx, http.StatusNotFound, "not found")
}

func failInternal(ctx *gin.Context, err error) {
	fail(ctx, http.StatusInternalServerError, err.Error())
}

func failBadRequest(ctx *gin.Context) {
	fail(ctx, http.StatusBadRequest, "bad request")
}

func failUnauthorized(ctx *gin.Context) {
	fail(ctx, http.StatusUnauthorized, "unauthorized")
}

func ok(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
	})
}

func (r *RestAPI) handlerCORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Header("Access-Control-Allow-Methods", "POST,GET,DELETE")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,Cookie")
	ctx.Header("Access-Control-Allow-Credentials", "true")
}

func (r *RestAPI) handlePreflight(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodOptions {
		ctx.Status(http.StatusOK)
		ctx.Abort()
		return
	}
	ctx.Next()
}

func checkFile(fileInfo *storage.File) error {
	if !util.StringArrayContains(allowedMimeTypes, fileInfo.ContentType) {
		return errInvalidMimeType
	}

	if !allowedFileName.MatchString(fileInfo.Name) {
		return errInvalidFileName
	}

	if fileInfo.Size > allowedFileSize {
		return errInvalidFileSize
	}

	return nil
}
