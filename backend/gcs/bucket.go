package gcs

import (
	"fmt"
	"os"
)

func PublicBucket() (string, error) {
	const varName = "PUBLIC_GCS_BUCKET"
	b := os.Getenv(varName)
	if b == "" {
		return "", fmt.Errorf("%s environment variable is unset", varName)
	}
	return b, nil
}
