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
	for _, export := range m.baseData.PerUserEntries {
		err := m.datastore.SetUserProfile(export.Username, types.UserProfile{
			AboutMarkdown: types.UserBio(export.Profile.About),
			EmailAddress:  types.EmailAddress(export.Profile.Email),
			TwitterHandle: types.TwitterHandle(export.Profile.Twitter),
		})
		if err != nil {
			return err
		}
		for _, d := range export.Drafts {
			err := m.datastore.InsertDraft(export.Username, d)
			if err != nil {
				return err
			}
		}
		for _, e := range export.Entries {
			err := m.datastore.InsertEntry(export.Username, e)
			if err != nil {
				return err
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
