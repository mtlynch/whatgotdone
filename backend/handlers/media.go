package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/dates"
	"github.com/mtlynch/whatgotdone/backend/gcs"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s *defaultServer) mediaPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.gcsClient == nil {
			log.Printf("can't accept media upload because media uploads are disabled")
			http.Error(w, "Media uploading is disabled", http.StatusBadRequest)
			return
		}

		mediaFile, contentType, err := mediaFileFromRequest(w, r)
		if err != nil {
			log.Printf("failed to read media from request: %v", err)
			http.Error(w, fmt.Sprintf("Media upload failed: %v", err), http.StatusBadRequest)
			return
		}

		username := usernameFromContext(r.Context())
		path, err := mediaPath(contentType, username)
		if err != nil {
			log.Printf("failed to generate media path: %v", err)
			http.Error(w, fmt.Sprintf("Media upload failed: %v", err), http.StatusInternalServerError)
			return
		}

		url, err := s.gcsClient.UploadFile(mediaFile, path, contentType, gcs.CacheControlPublic)
		if err != nil {
			log.Printf("failed to read media from request: %v", err)
			http.Error(w, fmt.Sprintf("Media upload failed: %v", err), http.StatusBadRequest)
			return
		}

		respondOK(w, struct {
			URL string `json:"url"`
		}{
			URL: url,
		})
	}
}

const maxMediaSize = 20971520 // 20 MB

func mediaFileFromRequest(w http.ResponseWriter, r *http.Request) (io.Reader, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxMediaSize)
	r.ParseMultipartForm(32 << 20)
	file, metadata, err := r.FormFile("file")
	if err != nil {
		return nil, "", err
	}
	contentType := metadata.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return nil, "", errors.New("invalid media format")
	}
	return file, contentType, nil
}

func mediaPath(contentType string, username types.Username) (string, error) {
	timestamp := dates.ThisFriday().Format("20060102")
	extension := "png"
	if contentType == "image/jpeg" {
		extension = "jpg"
	}

	return fmt.Sprintf("uploads/%s/%s/%s.%s", string(username), timestamp, randomKey(), extension), nil
}

func randomKey() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
