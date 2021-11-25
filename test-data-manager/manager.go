// +build dev staging

package main

import (
	"database/sql"
	"log"
	"os"

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
	log.Printf("%+v", m.baseData.PerUserEntries)
	wipeDb()
	for _, perUserEntries := range m.baseData.PerUserEntries {
		for _, d := range perUserEntries.Drafts {
			err := m.datastore.InsertDraft(perUserEntries.Username, d)
			if err != nil {
				return err
			}
		}
		for _, e := range perUserEntries.Entries {
			err := m.datastore.InsertEntry(perUserEntries.Username, e)
			if err != nil {
				return err
			}
		}
	}
	for u, p := range m.baseData.Profiles {
		err := m.datastore.SetUserProfile(types.Username(u), types.UserProfile{
			AboutMarkdown: types.UserBio(p.About),
			EmailAddress:  types.EmailAddress(p.Email),
			TwitterHandle: types.TwitterHandle(p.Twitter),
		})
		if err != nil {
			return err
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
