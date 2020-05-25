package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/zekroTJA/yuri2/internal/auth"
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/discord"
	"github.com/zekroTJA/yuri2/internal/restapi"
	"github.com/zekroTJA/yuri2/internal/storage"
	"github.com/zekroTJA/yuri2/pkg/discordoauth"
)

var (
	flagEnvFile = flag.String("env", "", ".env file location")
)

func main() {
	flag.Parse()

	cfg, err := config.Parse("YURI", *flagEnvFile)
	if err != nil {
		log.Fatalf("failed initializing config: %s", err.Error())
	}

	minioCfg := storage.MinioConfig{
		AccessKey: cfg.StorageKeyID,
		Secret:    cfg.StorageAccessKey,
		Endpoint:  cfg.StorageEndpoint,
		Location:  cfg.StorageLocation,
		UseSSL:    cfg.StorageUseSSL,
	}
	store := new(storage.Minio)
	if err = store.Init(&minioCfg); err != nil {
		log.Fatalf("failed initializing storage: %s", err.Error())
	}

	db, err := database.NewPostgres(cfg.DatabaseConnStr)
	if err != nil {
		log.Fatalf("failed initializing database: %s", err.Error())
	}

	dg := discord.New(cfg.DiscordToken, cfg.DiscordPrefix, db)
	go dg.RunBlocking()

	doa := discordoauth.New(
		cfg.DiscordID, cfg.DiscordSecret,
		fmt.Sprintf("%s/%s", cfg.WsPublicAddr, "auth/callback"),
		"identify")

	aut, err := auth.NewJWTAuth(cfg.WsAuthSignKey, &auth.Options{
		ExpireTime: 30 * time.Hour * 24,
	})
	if err != nil {
		log.Fatalf("failed initializing auth module: %s", err.Error())
	}

	rapi, err := restapi.New(store, db, aut, doa, dg)
	if err != nil {
		log.Fatalf("failed initializing webserver: %s", err.Error())
	}

	err = rapi.ListenAndServeBlocking(cfg.WsAddr)
	if err != nil {
		log.Fatalf("failed binding webserver: %s", err.Error())
	}
}
