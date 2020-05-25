package discord

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/zekroTJA/yuri2/internal/database"
)

type Discord struct {
	db database.Database

	handler *cmdHandler
	client  *disgord.Client
	myself  *disgord.User
}

func New(token, prefix string, db database.Database) (d *Discord) {
	d = new(Discord)
	d.db = db
	d.handler = newCmdHandler(d, prefix)
	d.client = disgord.New(disgord.Config{
		BotToken: token,
	})

	d.registerHandlers()

	return d
}

func (d *Discord) registerHandlers() {
	d.client.On(disgord.EvtReady, d.hReady)
	d.client.On(disgord.EvtMessageCreate, d.handler.Handler)
}

func (d *Discord) RunBlocking() error {
	return d.client.StayConnectedUntilInterrupted(context.Background())
}
