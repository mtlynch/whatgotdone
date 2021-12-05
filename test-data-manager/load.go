//go:build dev || staging

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/mtlynch/whatgotdone/backend/types/export"
)

type (
	initData struct {
		UserData map[string]export.UserData `json:"userData" yaml:"user_data"`
	}
)

func loadFromFile(filename string) initData {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var d initData

	unmarshal := yaml.Unmarshal
	if strings.HasSuffix(filename, ".json") {
		unmarshal = json.Unmarshal
	}

	err = unmarshal(b, &d)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return d
}
