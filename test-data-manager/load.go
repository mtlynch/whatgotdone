// +build dev staging

package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/mtlynch/whatgotdone/backend/types"
)

type (
	userEntries struct {
		Username types.Username       `yaml:"username"`
		Drafts   []types.JournalEntry `yaml:"drafts"`
		Entries  []types.JournalEntry `yaml:"entries"`
	}

	initData struct {
		PerUserEntries []userEntries `yaml:"perUserEntries"`
	}
)

func loadYaml(filename string) initData {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var d initData
	err = yaml.Unmarshal(b, &d)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return d
}
