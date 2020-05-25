package restapi

import (
	"crypto/rand"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zekroTJA/yuri2/internal/auth"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/discord"
	"github.com/zekroTJA/yuri2/internal/storage"
	"github.com/zekroTJA/yuri2/pkg/discordoauth"
)

type RestAPI struct {
	s    storage.Storage
	db   database.Database
	auth auth.Auth
	doa  *discordoauth.DiscordOauth
	dg   *discord.Discord

	e   *gin.Engine
	rak string
}

func New(store storage.Storage, db database.Database, aut auth.Auth, doa *discordoauth.DiscordOauth, dg *discord.Discord) (r *RestAPI, err error) {
	r = new(RestAPI)
	r.e = gin.Default()
	r.s = store
	r.db = db
	r.auth = aut
	r.doa = doa
	r.dg = dg

	r.e.MaxMultipartMemory = allowedFileSize

	r.registerHandlers()
	if err = r.generateRAK(); err != nil {
		return
	}

	// TODO: Remove in release build
	fmt.Println("RAK: ", r.rak)

	return
}

func (r *RestAPI) ListenAndServeBlocking(addr string) error {
	return r.e.Run(addr)
}

func (r *RestAPI) registerHandlers() {
	r.e.Group("/auth").
		GET("/login", r.hAuthLogin).
		GET("/callback", r.hAuthCallback)

	r.e.Group("/sounds", r.hAuthValidate).
		GET("", r.hSoundsGet).
		POST("", r.hSoundsPost).
		GET("/:fileName", r.hSoundGet)
}

func (r *RestAPI) generateRAK() error {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return err
	}

	r.rak = fmt.Sprintf("%x", key)

	return nil
}
