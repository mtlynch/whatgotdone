package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
)

// InsertPageViews stores the set of pageview data for What Got Done routes.
func (d db) InsertPageViews(pvcs []ga.PageViewCount) error {
	log.Printf("saving %d page view entries to datastore", len(pvcs))
	valueClauses := strings.TrimSuffix(strings.Repeat("(?,?,strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc')), ", len(pvcs)), ", ")
	values := make([]interface{}, len(pvcs)*2)
	for i, pvc := range pvcs {
		values[i*2] = pvc.Path
		values[(i*2)+1] = pvc.Views
	}
	_, err := d.ctx.Exec(fmt.Sprintf(`
	INSERT OR REPLACE INTO pageviews(
		path,
		views,
		last_updated
	)
	VALUES
		%s`, valueClauses), values...)
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
