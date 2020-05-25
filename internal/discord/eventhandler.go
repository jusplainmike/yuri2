package discord

import (
	"log"

	"github.com/andersfylling/disgord"
)

func (d *Discord) hReady(s disgord.Session, e *disgord.Ready) {
	var err error
	if d.myself, err = s.GetCurrentUser(e.Ctx); err != nil {
		log.Println("ERR :", err)
	}

	log.Println("INFO : READY")
}
