package gcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func (c Client) DeletePath(path string) error {
	log.Printf("Deleting gs://%s/%s", c.bucketName, path)
	ctx := context.Background()
	bh := c.gcsClient.Bucket(c.bucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	it := bh.Objects(ctx, &storage.Query{
		Prefix: path,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate objects under %s: %v", path, err)
		}
		obj := bh.Object(attrs.Name)
		if err := obj.Delete(ctx); err != nil {
			return fmt.Errorf("failed to delete %s: %v", attrs.Name, err)
		}
	}

	return nil
}
