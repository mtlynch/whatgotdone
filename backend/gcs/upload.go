package gcs

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

// CacheControlPublic indicates that any cache can store the response.
const CacheControlPublic = "public"

// UploadFile uploads a file to a lcoation on Google Cloud Storage.
func (c Client) UploadFile(r io.Reader, path, contentType, cacheControl string) (string, error) {
	log.Printf("Saving image to gs://%s/%s", c.bucketName, path)
	ctx := context.Background()
	bh := c.gcsClient.Bucket(c.bucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	obj := bh.Object(path)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	_, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{
		ContentType:  contentType,
		CacheControl: cacheControl,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, path), nil
}
