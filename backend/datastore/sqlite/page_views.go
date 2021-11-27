package sqlite

import (
	"database/sql"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

// InsertPageViews stores the count of pageviews for a given What Got Done route.
func (d db) InsertPageViews(path string, pageViews int) error {
	_, err := d.ctx.Exec(`
	INSERT OR REPLACE INTO pageviews(
		path,
		views,
		last_updated)
	values(?,?,strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc'))`, path, pageViews)
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
		return 0, datastore.PageViewsNotFoundError{Path: path}
	} else if err != nil {
		return 0, err
	}

	return pageViews, nil
}
