package firestore

// Close cleans up datastore resources. Clients should not call any Datastore
// functions after calling Close().
func (c client) Close() error {
	return c.firestoreClient.Close()
}
