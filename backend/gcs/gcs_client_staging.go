// +build staging

package gcs

import (
	"context"
	"errors"

	"cloud.google.com/go/storage"
)

func newGcsClient(_ context.Context) (*storage.Client, error) {
	return nil, errors.New("staging environment does not support GCS client")
}
