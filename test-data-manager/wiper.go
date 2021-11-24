// +build dev staging

package main

import (
	"log"
	"os"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type wiper struct {
}

func newWiper() wiper {
	return wiper{}
}

func (w *wiper) Wipe() {
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
