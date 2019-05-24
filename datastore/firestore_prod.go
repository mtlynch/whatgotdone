// +build !dev
// +build !staging

package datastore

import (
	"context"

	"cloud.google.com/go/firestore"
)

func newFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	const firestoreProjectID = "whatgotdone"
	return firestore.NewClient(ctx, firestoreProjectID)
}
