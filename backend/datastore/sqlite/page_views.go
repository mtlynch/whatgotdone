package sqlite

import (
	"database/sql"
	"log"
)

// InsertPageViews stores the count of pageviews for a given What Got Done route.
func (d db) InsertPageViews(path string, pageViews int) error {
	log.Printf("saving pageviews to datastore: %s : %d", path, pageViews)
	_, err := d.ctx.Exec(`
	INSERT OR REPLACE INTO pageviews(
		path,
		views,
		last_updated)
	values(?,?,datetime('now'))`, path, pageViews)
	return err
}

// GetPageViews retrieves the count of pageviews for a given What Got Done route.
func (d db) GetPageViews(path string) (int, error) {
	stmt, err := d.ctx.Prepare(`
	SELECT
		views
	FROM
		pageviews
	WHERE
		path=?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var pageViews int
	err = stmt.QueryRow(path).Scan(&pageViews)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return pageViews, nil
}
