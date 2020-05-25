package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WsAddr        string `envconfig:"WEBSERVER_ADDRESS" default:"localhost:8080"`
	WsPublicAddr  string `envconfig:"WEBSERVER_PUBLIC_ADDRESS" default:"http://localhost:8080"`
	WsAuthSignKey []byte `envconfig:"WEBSERVER_AUTH_SIGN_KEY"`

	StorageEndpoint  string `envconfig:"STORAGE_ENDPOINT" required:"true"`
	StorageKeyID     string `envconfig:"STORAGE_KEY_ID" required:"true"`
	StorageAccessKey string `envconfig:"STORAGE_ACCESS_KEY" required:"true"`
	StorageLocation  string `envconfig:"STORAGE_LOCATION" required:"true"`
	StorageUseSSL    bool   `envconfig:"STORAGE_USE_SSL"`

	DiscordToken  string `envconfig:"DISCORD_TOKEN" required:"true"`
	DiscordID     string `envconfig:"DISCORD_ID" required:"true"`
	DiscordSecret string `envconfig:"DISCORD_SECRET" required:"true"`
	DiscordPrefix string `envconfig:"DISCORD_PREFIX" default:"y!"`

	DatabaseConnStr string `envconfig:"DATABASE_CONNECTION_STRING" required:"true"`
}

func Parse(prefix string, envFiles ...string) (cfg *Config, err error) {
	godotenv.Load(envFiles...)
	cfg = new(Config)
	err = envconfig.Process(prefix, cfg)
	return
}
