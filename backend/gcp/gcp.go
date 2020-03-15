package gcp

import (
	"log"
	"os"
)

func ProjectID() string {
	return mustGetenv("GOOGLE_CLOUD_PROJECT")
}

func mustGetenv(key string) string {
	e := os.Getenv(key)
	if e == "" {
		log.Fatalf("%s environment variable must be set", key)
	}
	return e
}
