package gcs

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/mtlynch/whatgotdone/backend/gcp"
)

func newGcsClient(ctx context.Context) (*storage.Client, error) {
	log.Printf("loading GCS client with service account: %s", gcp.ServiceAccountKeyFile)
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcp.ServiceAccountKeyFile))
	if err != nil {
		return nil, err
	}

	return client, nil
}
