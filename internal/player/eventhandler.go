package player

import (
	"github.com/foxbot/gavalink"
)

type EventHandler interface {
	OnTrackStart(player *gavalink.Player, track, ident string, resource ResourceType, guildID, userID, userTag string)
	OnTrackEnd(player *gavalink.Player, track string, reason string) error
	OnTrackException(player *gavalink.Player, track string, reason string) error
	OnTrackStuck(player *gavalink.Player, track string, threshold int) error
	OnVolumeChanged(player *gavalink.Player, guildID string, vol int)
}

type EventHandlerManager struct {
	handler []EventHandler
}

func NewEventHandlerManager() *EventHandlerManager {
	return &EventHandlerManager{
		handler: make([]EventHandler, 0),
	}
}

func (h *EventHandlerManager) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	for _, hnd := range h.handler {
		hnd.OnTrackEnd(player, track, reason)
	}

	return nil
}

func (h *EventHandlerManager) OnTrackException(player *gavalink.Player, track string, reason string) error {
	for _, hnd := range h.handler {
		hnd.OnTrackException(player, track, reason)
	}

	return nil
}

func (h *EventHandlerManager) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	for _, hnd := range h.handler {
		hnd.OnTrackStuck(player, track, threshold)
	}

	return nil
}

func (h *EventHandlerManager) OnTrackStart(player *gavalink.Player, track, ident string,
	resource ResourceType, guildID, userID, userTag string) {

	for _, hnd := range h.handler {
		hnd.OnTrackStart(player, track, ident, resource, guildID, userID, userTag)
	}
}

func (h *EventHandlerManager) OnVolumeChanged(player *gavalink.Player, guildID string, vol int) {
	for _, hnd := range h.handler {
		hnd.OnVolumeChanged(player, guildID, vol)
	}
}

func (h *EventHandlerManager) AddHandler(handler EventHandler) {
	h.handler = append(h.handler, handler)
}
