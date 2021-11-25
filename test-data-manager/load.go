// +build dev staging

package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/mtlynch/whatgotdone/backend/types/export"
)

type (
	initData struct {
		UserData map[string]export.UserData `yaml:"user_data"`
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
