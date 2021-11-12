package sqlite

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

type db struct {
	ctx *sql.DB
}

var notImplementedError = errors.New("not implemented")

func New() datastore.Datastore {
	dbDir := "data"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}
	ctx, err := sql.Open("sqlite3", dbDir+"/store.db")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = ctx.Exec(`
	CREATE TABLE IF NOT EXISTS user_preferences (
		username TEXT PRIMARY KEY,
		entry_template TEXT
		)
	`)
	if err != nil {
		log.Fatalln(err)
	}
	return &db{
		ctx: ctx,
	}
}
