package gcs

import (
	"context"

	"cloud.google.com/go/storage"
)

type Client struct {
	gcsClient  *storage.Client
	bucketName string
}

func New() (*Client, error) {
	b, err := PublicBucket()
	if err != nil {
		return &Client{}, err
	}

	client, err := newGcsClient(context.Background())
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		gcsClient:  client,
		bucketName: b,
	}, nil
}
