package discord

import (
	"github.com/andersfylling/disgord"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/yuri2/internal/static"
)

type embed struct {
	emb *disgord.Embed

	msg *disgord.Message
}

func Embed(content string) *embed {
	e := new(embed)
	e.emb = &disgord.Embed{
		Color:       static.ColorMain,
		Description: content,
	}

	return e
}

func (e *embed) Send(channel snowflake.ID) {

}
