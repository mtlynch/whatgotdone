// +build !dev
// +build !staging

package client

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"

	"github.com/mtlynch/whatgotdone/backend/gcp"
)

func New(ctx context.Context) (*firestore.Client, error) {
	log.Printf("creating production-mode firestore client")
	return firestore.NewClient(ctx, gcp.ProjectID())
}
