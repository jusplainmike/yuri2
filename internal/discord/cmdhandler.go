package discord

import (
	"strings"

	"github.com/andersfylling/disgord"
)

type Command interface {
	Invokes() []string
	Description() string
	Help() string
	Exec(d *Discord, m *disgord.Message, args []string) error
}

type cmdHandler struct {
	d      *Discord
	prefix string

	cmds []Command
}

func newCmdHandler(d *Discord, prefix string) *cmdHandler {
	return &cmdHandler{
		prefix: prefix,
		d:      d,
	}
}

func (c *cmdHandler) Register(cmd ...Command) {
	c.cmds = append(c.cmds, cmd...)
}

func (c *cmdHandler) Handler(s disgord.Session, e *disgord.MessageCreate) {
	if c.d.myself.ID == e.Message.Author.ID || e.Message.Author.Bot {
		return
	}

	prefix, _ := c.d.db.GetPrefix(e.Message.GuildID.String())
	if prefix == "" {
		prefix = c.prefix
	}

	if !strings.HasPrefix(e.Message.Content, prefix) {
		return
	}

	split := strings.Split(e.Message.Content, " ")
	invoke := split[0]

	if len(invoke) <= len(prefix) {
		return
	}

	invoke = invoke[len(prefix):]

	if cmd, ok := c.getCmdByInvoke(invoke); ok {
		cmd.Exec(c.d, e.Message, split[1:])
	}
}

func (c *cmdHandler) getCmdByInvoke(invoke string) (Command, bool) {
	invoke = strings.ToLower(invoke)

	for _, cmd := range c.cmds {
		for _, i := range cmd.Invokes() {
			if i == invoke {
				return cmd, true
			}
		}
	}

	return nil, false
}
