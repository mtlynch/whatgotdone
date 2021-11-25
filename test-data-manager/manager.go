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
	log.Printf("%+v", m.baseData.PerUserEntries)
	m.wiper.Wipe()
	for _, perUserEntries := range m.baseData.PerUserEntries {
		for _, d := range perUserEntries.Drafts {
			err := m.datastore.InsertDraft(perUserEntries.Username, types.JournalEntry{
				Date:     d.Date,
				Markdown: d.Markdown,
			})
			if err != nil {
				return err
			}
		}
		for _, e := range perUserEntries.Entries {
			err := m.datastore.InsertEntry(perUserEntries.Username, types.JournalEntry{
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
