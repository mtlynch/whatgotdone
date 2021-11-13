package sqlite

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

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
		);
	CREATE TABLE IF NOT EXISTS user_profiles (
		username TEXT PRIMARY KEY,
		about_markdown TEXT,
		email TEXT,
		twitter TEXT,
		mastodon TEXT
		);
	CREATE TABLE IF NOT EXISTS journal_entries(
		username TEXT,
		date TEXT,
		last_modified TEXT,
		markdown TEXT,
		is_draft INTEGER,
		PRIMARY KEY (username, date, is_draft)
		);
	`)
	if err != nil {
		log.Fatalln(err)
	}
	return &db{
		ctx: ctx,
	}
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func parseDatetime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}
