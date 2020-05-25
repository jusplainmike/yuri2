package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connStr string) (p *Postgres, err error) {
	p = new(Postgres)
	p.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return
	}

	return p, p.setupStructure()
}

func (p *Postgres) setupStructure() (err error) {
	// GuildProps Table
	if _, err = p.db.Exec(
		`CREATE TABLE IF NOT EXISTS public."GuildProps" (
			"GuildId" varchar(25) NOT NULL,
			"UploaderRoleId" varchar(25),
			"AdminRoleId" varchar(25),
			"MuteRoleId" varchar(25),
			"Prefix" varchar(10),
			PRIMARY KEY ("GuildId") 
		 );`,
	); err != nil {
		return
	}

	// UserProps Table
	if _, err = p.db.Exec(
		`CREATE TABLE IF NOT EXISTS public."UserProps" (
			"UserId" varchar(25) NOT NULL,
			"FastTrigger" varchar(60),
			PRIMARY KEY ("UserId") 
		);`,
	); err != nil {
		return
	}

	// SoundLog Table
	if _, err = p.db.Exec(
		`CREATE TABLE IF NOT EXISTS public."SoundLog" (
			"Id" SERIAL NOT NULL,
			"GuildId" varchar(25),
			"GuildName" varchar(100),
			"ExecutorId" varchar(25),
			"ExecutorName" varchar(100),
			"Sound" varchar(60),
			"Type" smallint NOT NULL DEFAULT '0',
			"TimeStamp" timestamp,
			PRIMARY KEY ("Id") 
		);`,
	); err != nil {
		return
	}

	return
}

func (p *Postgres) GetUploaderRole(guildId string) (string, error) {
	var val string
	err := p.db.QueryRow(
		`SELECT "UploaderRoleId" FROM public."GuildProps" WHERE "GuildId" = $1;`, guildId).
		Scan(&val)
	return val, err
}

func (p *Postgres) SetUploaderRole(guildId, roleId string) error {
	_, err := p.db.Exec(
		`UPDATE public."GuildProps" SET "UploaderRoleId" = $2 WHERE "GuildId" = $1;`, guildId, roleId)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(
		`INSERT INTO public."GuildProps" ("GuildId", "UploaderRoleId") 
			SELECT $1, $2 
			WHERE NOT EXISTS (SELECT 1 FROM public."GuildProps" WHERE "GuildId" = $3)`,
		guildId, roleId, guildId)
	return err
}

func (p *Postgres) GetAdminRole(guildId string) (string, error) {
	var val string
	err := p.db.QueryRow(
		`SELECT "AdminRoleId" FROM public."GuildProps" WHERE "GuildId" = $1;`, guildId).
		Scan(&val)
	return val, err
}

func (p *Postgres) SetAdminRole(guildId, roleId string) error {
	_, err := p.db.Exec(
		`UPDATE public."GuildProps" SET "AdminRoleId" = $2 WHERE "GuildId" = $1;`, guildId, roleId)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(
		`INSERT INTO public."GuildProps" ("GuildId", "AdminRoleId") 
			SELECT $1, $2
			WHERE NOT EXISTS (SELECT 1 FROM public."GuildProps" WHERE "GuildId" = $3)`, guildId, roleId, guildId)
	return err
}

func (p *Postgres) GetMutedRole(guildId string) (string, error) {
	var val string
	err := p.db.QueryRow(
		`SELECT "MuteRoleId" FROM public."GuildProps" WHERE "GuildId" = $1;`, guildId).
		Scan(&val)
	return val, err
}

func (p *Postgres) SetMutedRole(guildId, roleId string) error {
	_, err := p.db.Exec(
		`UPDATE public."GuildProps" SET "MuteRoleId" = $2 WHERE "GuildId" = $1;`, guildId, roleId)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(
		`INSERT INTO public."GuildProps" ("GuildId", "MuteRoleId") 
			SELECT $1, $2 
			WHERE NOT EXISTS (SELECT 1 FROM public."GuildProps" WHERE "GuildId" = $3)`, guildId, roleId, guildId)
	return err
}

func (p *Postgres) GetPrefix(guildId string) (string, error) {
	var val string
	err := p.db.QueryRow(
		`SELECT "Prefix" FROM public."GuildProps" WHERE "GuildId" = $1;`, guildId).
		Scan(&val)
	return val, err
}

func (p *Postgres) SetPrefix(guildId, prefix string) error {
	_, err := p.db.Exec(
		`UPDATE public."GuildProps" SET "Prefix" = $2 WHERE "GuildId" = $1;`, guildId, prefix)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(
		`INSERT INTO public."GuildProps" ("GuildId", "Prefix") 
			SELECT $1, $2
			WHERE NOT EXISTS (SELECT 1 FROM public."GuildProps" WHERE "GuildId" = $3)`, guildId, prefix, guildId)
	return err
}

func (p *Postgres) GetShortTrigger(userId string) (string, error) {
	var val string
	err := p.db.QueryRow(
		`SELECT "FastTrigger" FROM public."UserProps" WHERE "UserId" = $1;`, userId).
		Scan(&val)
	return val, err
}

func (p *Postgres) SetShortTrigger(userId, sound string) error {
	_, err := p.db.Exec(
		`UPDATE public."UserProps" SET "FastTrigger" = $2 WHERE "UserId" = $1;`, userId, sound)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(
		`INSERT INTO public."UserProps" ("UserId", "FastTrigger") 
			SELECT $1, $2 
			WHERE NOT EXISTS (SELECT 1 FROM public."UserProps" WHERE "UserId" = $3)`, userId, sound, userId)
	return err
}

func (p *Postgres) SoundPlayed(entry *SoundLogEntry) error {
	_, err := p.db.Exec(
		`INSERT INTO public."SoundLog" (
			"GuildId", "GuildName", "ExecutorId", "ExecutorName",
			"Sound", "Type", "TimeStamp"
		) VALUES ( $1, $2, $3, $4, $5, $6, $7 );`,
		entry.GuildId, entry.GuildName, entry.ExecutorId, entry.ExecutorName,
		entry.Sound, entry.Type, entry.TimeStamp,
	)

	return err
}

func (p *Postgres) GetSoundLog(guildId string, limit, offset int) ([]*SoundLogEntry, error) {
	out := make([]*SoundLogEntry, limit)
	var i int

	rows, err := p.db.Query(
		`SELECT 
			"GuildId", "GuildName", "ExecutorId", "ExecutorName",
			"Sound", "Type", "TimeStamp"
		FROM public."SoundLog"
		WHERE "GuildId" = $1 
		ORDER BY "Id" DESC OFFSET $2 LIMIT $3;`,
		guildId, offset, limit,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		entry := new(SoundLogEntry)
		err = rows.Scan(&entry.GuildId, &entry.GuildName, &entry.ExecutorId, &entry.ExecutorName,
			&entry.Sound, &entry.Type, &entry.TimeStamp)
		if err != nil {
			return nil, err
		}
		out[i] = entry
		i++
	}

	return out[0:i], nil
}

func (p *Postgres) GetSoundLogLen(guildId string) (int, error) {
	var c int
	err := p.db.QueryRow(
		`SELECT COUNT(*) FROM public."SoundLog"
		WHERE "GuildId" = $1;`, guildId,
	).Scan(&c)

	return c, err
}

func (p *Postgres) GetSoundStats(guildId string, limit int) ([]*SoundStatsEntry, error) {
	ln, err := p.GetSoundLogLen(guildId)
	if err != nil {
		return nil, err
	}

	logEntries, err := p.GetSoundLog(guildId, ln, 0)
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]*SoundStatsEntry)
	for _, e := range logEntries {
		c, ok := countMap[e.Sound]
		if !ok {
			c = &SoundStatsEntry{
				GuildId:   e.GuildId,
				GuildName: e.GuildName,
				N:         0,
				Sound:     e.Sound,
			}
			countMap[e.Sound] = c
		}
		c.N++
		if e.TimeStamp.After(c.LastPlayed) {
			c.LastPlayed = e.TimeStamp
		}
	}

	lCountMap := len(countMap)
	if lCountMap < limit {
		limit = lCountMap
	}

	out := make([]*SoundStatsEntry, limit)
	var i int
	for _, v := range countMap {
		out[i] = v
		i++
	}

	return out, nil
}
