// +build staging

package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func newFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	const devServiceAccount = "service-account-creds-staging.json"
	return firestore.NewClient(ctx, getGoogleCloudProjectId(), option.WithCredentialsFile(devServiceAccount))
}
