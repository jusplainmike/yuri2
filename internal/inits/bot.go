package inits

import (
	"github.com/zekroTJA/discordgocmds"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/discordbot/commands"
	"github.com/zekroTJA/yuri2/internal/logger"
)

// InitDiscordBot initializes the discord bot session and registers
// all commands and handlers.
func InitDiscordBot(token, ownerID, generalPrefix string, dbMiddleware database.Middleware) *discordbot.Bot {
	handlers := []interface{}{}

	commands := []discordgocmds.Command{
		&commands.Prefix{PermLvl: 5, DB: dbMiddleware},
	}

	bot, err := discordbot.NewBot(token, ownerID, generalPrefix, dbMiddleware)
	if err != nil {
		logger.Fatal("DBOT :: failed initialization: %s", err.Error())
	}

	bot.RegisterHandler(handlers)
	bot.RegisterCommands(commands)

	logger.Info("DBOT :: initialized")

	if err := bot.Open(); err != nil {
		logger.Fatal("DBOT :: failed connecting to the discord API: %s", err.Error())
	}

	logger.Info("DBOT :: connection established")

	return bot
}