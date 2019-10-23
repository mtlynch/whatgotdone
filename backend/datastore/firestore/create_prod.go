// +build !dev
// +build !staging

package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

func newFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	return firestore.NewClient(ctx, getGoogleCloudProjectID())
}
