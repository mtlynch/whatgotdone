// +build dev

package datastore

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func newFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	const firestoreProjectID = "whatgotdone-dev"
	const devServiceAccount = "service-account-creds-dev.json"
	return firestore.NewClient(ctx, firestoreProjectID, option.WithCredentialsFile(devServiceAccount))
}
