//go:build dev || staging

package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/datastore/sqlite"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type manager struct {
	datastore datastore.Datastore
	baseData  initData
}

func NewManager(baseData initData) manager {
	return manager{
		datastore: sqlite.New(),
		baseData:  baseData,
	}
}

func (m *manager) Reset() error {
	log.Printf("resetting datastore data")
	wipeDb()
	for username, ud := range m.baseData.UserData {
		err := m.datastore.SetPreferences(types.Username(username), types.Preferences{
			EntryTemplate: types.EntryContent(ud.Preferences.EntryTemplate),
		})
		if err != nil {
			return err
		}
		err = m.datastore.SetUserProfile(types.Username(username), types.UserProfile{
			AboutMarkdown:   ud.Profile.AboutMarkdown,
			EmailAddress:    ud.Profile.EmailAddress,
			TwitterHandle:   ud.Profile.TwitterHandle,
			MastodonAddress: ud.Profile.MastodonAddress,
		})
		if err != nil {
			return err
		}
		for _, d := range ud.Drafts {
			err := m.datastore.InsertDraft(types.Username(username), types.JournalEntry{
				Date:     d.Date,
				Markdown: types.EntryContent(d.Markdown),
			})
			if err != nil {
				return err
			}
		}
		for _, e := range ud.Entries {
			err := m.datastore.InsertEntry(types.Username(username), types.JournalEntry{
				Date:     e.Date,
				Markdown: types.EntryContent(e.Markdown),
			})
			if err != nil {
				return err
			}
		}
		for _, leader := range ud.Following {
			err := m.datastore.InsertFollow(leader, types.Username(username))
			if err != nil {
				return err
			}
		}
		for date, reactions := range ud.Reactions {
			for _, r := range reactions {
				ts, err := time.ParseInLocation(time.RFC3339, r.Timestamp, time.UTC)
				if err != nil {
					return err
				}
				err = m.datastore.AddReaction(types.Username(username), date, types.Reaction{
					Username:  r.Username,
					Symbol:    r.Symbol,
					Timestamp: ts,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func wipeDb() {
	dbDir := "data"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}
	ctx, err := sql.Open("sqlite3", dbDir+"/store.db")
	if err != nil {
		log.Fatalln(err)
	}
	tables := []string{
		"user_preferences",
		"user_profiles",
		"journal_entries",
		"follows",
		"entry_reactions",
		"pageviews",
	}
	for _, tbl := range tables {
		_, err = ctx.Exec("DELETE FROM " + tbl)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
