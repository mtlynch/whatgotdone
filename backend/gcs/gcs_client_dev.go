// +build dev

package gcs

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func newGcsClient(ctx context.Context) (*storage.Client, error) {
	const devServiceAccount = "service-account-creds-dev.json"
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(devServiceAccount))
	if err != nil {
		return nil, err
	}
	log.Print("loaded dev GCS credentials")

	return client, nil
}
