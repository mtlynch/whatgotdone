package datastore

import (
	"context"
	"log"
)

type (
	UserKitKeystore interface {
		SecretKey() (string, error)
	}

	secretKey struct {
		Value string `firestore:"value,omitempty"`
	}
)

func NewUserKitKeyStore() UserKitKeystore {
	ctx := context.Background()

	client, err := newFirestoreClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return &defaultClient{
		firestoreClient: client,
		ctx:             ctx,
	}
}
func (c defaultClient) SecretKey() (string, error) {
	doc, err := c.firestoreClient.Collection(secretsRootKey).Doc(secretUserKitDocKey).Get(c.ctx)
	if err != nil {
		return "", err
	}
	var sk secretKey
	doc.DataTo(&sk)
	return sk.Value, nil
}
