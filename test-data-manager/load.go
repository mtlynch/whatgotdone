// +build dev staging

package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/mtlynch/whatgotdone/backend/types"
)

type (
	profile struct {
		About   string `yaml:"about"`
		Email   string `yaml:"email"`
		Twitter string `yaml:"twitter"`
	}

	userEntries struct {
		Username types.Username       `yaml:"username"`
		Profile  profile              `yaml:"profile"`
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
