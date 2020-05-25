package database

import "time"

const (
	PlayTypeStorage PlayType = iota
	PlayTypeYoutube
)

type PlayType int

type SoundLogEntry struct {
	GuildId      string    `json:"guildid"`
	GuildName    string    `json:"guildname"`
	ExecutorId   string    `json:"executorid"`
	ExecutorName string    `json:"executorname"`
	Sound        string    `json:"sound"`
	Type         PlayType  `json:"type"`
	TimeStamp    time.Time `json:"timestamp"`
}

type SoundStatsEntry struct {
	GuildId    string    `json:"guildid"`
	GuildName  string    `json:"guildname"`
	N          int       `json:"n"`
	Sound      string    `json:"sound"`
	LastPlayed time.Time `json:"lastplayed"`
}

type Database interface {
	GetUploaderRole(guildId string) (string, error)
	SetUploaderRole(guildId, roleId string) error

	GetAdminRole(guildId string) (string, error)
	SetAdminRole(guildId, roleId string) error

	GetMutedRole(guildId string) (string, error)
	SetMutedRole(guildId, roleId string) error

	GetPrefix(guildId string) (string, error)
	SetPrefix(guildId, prefix string) error

	GetShortTrigger(userId string) (string, error)
	SetShortTrigger(userId, sound string) error

	SoundPlayed(entry *SoundLogEntry) error
	GetSoundLog(guildId string, limit, offset int) ([]*SoundLogEntry, error)
	GetSoundLogLen(guildId string) (int, error)
	GetSoundStats(guildId string, limit int) ([]*SoundStatsEntry, error)
}
