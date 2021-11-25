// +build dev staging

package main

import (
	"context"
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/datastore/firestore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func NewManager(baseData initData) manager {
	ctx := context.Background()
	return manager{
		datastore: firestore.New(),
		wiper:     newWiper(ctx),
		baseData:  baseData,
	}
}

type manager struct {
	datastore datastore.Datastore
	wiper     wiper
	baseData  initData
}

func (m *manager) Reset() error {
	log.Printf("resetting datastore data")
	m.wiper.Wipe()
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
				Date:     d.Date,
				Markdown: d.Markdown,
			})
			if err != nil {
				return err
			}
		}
		for _, e := range ud.Entries {
			err := m.datastore.InsertEntry(types.Username(username), types.JournalEntry{
				Date:     e.Date,
				Markdown: e.Markdown,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
