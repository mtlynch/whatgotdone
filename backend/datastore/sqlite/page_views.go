package sqlite

// InsertPageViews stores the count of pageviews for a given What Got Done route.
func (d db) InsertPageViews(path string, pageViews int) error {
	return notImplementedError
}

// GetPageViews retrieves the count of pageviews for a given What Got Done route.
func (d db) GetPageViews(path string) (int, error) {
	return 0, notImplementedError
}
