package sqlite

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

type DB struct {
	ctx *sql.DB
}

func New(path string) datastore.Datastore {
	log.Printf("reading DB from %s", path)
	ctx, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalln(err)
	}

	// The Litestream documentation recommends these pragmas.
	// https://litestream.io/tips/
	if _, err := ctx.Exec(`
		PRAGMA busy_timeout = 5000;
		PRAGMA synchronous = NORMAL;
		`); err != nil {
		log.Fatalf("failed to set pragmas: %v", err)
	}

	d := &DB{ctx: ctx}

	d.applyMigrations()

	return d
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func parseDatetime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05Z", s)
}
