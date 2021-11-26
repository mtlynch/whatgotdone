// +build dev staging

package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
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
			EntryTemplate: ud.Preferences.EntryTemplate,
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
				Date:         d.Date,
				Markdown:     d.Markdown,
				LastModified: canonicalizeDatetime(d.LastModified),
			})
			if err != nil {
				return err
			}
		}
		for _, e := range ud.Entries {
			err := m.datastore.InsertEntry(types.Username(username), types.JournalEntry{
				Date:         e.Date,
				Markdown:     e.Markdown,
				LastModified: canonicalizeDatetime(e.LastModified),
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
				err := m.datastore.AddReaction(types.Username(username), date, types.Reaction{
					Username:  r.Username,
					Symbol:    r.Symbol,
					Timestamp: r.Timestamp,
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

func canonicalizeDatetime(s string) string {
	t, err := parseDatetime(s)
	if err != nil {
		panic(err)
	}
	return t.Format("2006-01-02 15:04:05Z")
}

func parseDatetime(s string) (time.Time, error) {
	if strings.HasSuffix(s, "Z") {
		return time.ParseInLocation("2006-01-02T15:04:05Z", s, time.UTC)
	}
	return time.ParseInLocation("2006-01-02T15:04:05-07:00", s, time.UTC)
}
