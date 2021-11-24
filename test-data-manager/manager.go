// +build dev staging

package main

import (
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/datastore/sqlite"
)

func NewManager(baseData initData) manager {
	return manager{
		datastore: sqlite.New(),
		wiper:     newWiper(),
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
	return nil
}
