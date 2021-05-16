package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/dates"
)

func (s *defaultServer) mediaPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.gcsClient == nil {
			log.Printf("can't accept media upload because media uploads are disabled")
			http.Error(w, fmt.Sprintf("Media uploading is disabled"), http.StatusBadRequest)
			return
		}

		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must be logged in to upload an image or video", http.StatusForbidden)
			return
		}

		mediaFile, contentType, err := mediaFileFromRequest(w, r)
		if err != nil {
			log.Printf("failed to read media from request: %v", err)
			http.Error(w, fmt.Sprintf("Media upload failed: %v", err), http.StatusBadRequest)
			return
		}

		path, err := mediaPath(contentType, username)
		if err != nil {
			log.Printf("failed to generate media path: %v", err)
			http.Error(w, fmt.Sprintf("Media upload failed: %v", err), http.StatusInternalServerError)
			return
		}

		url, err := s.gcsClient.UploadFile(mediaFile, path, contentType, "public")
		if err != nil {
			log.Printf("failed to read media from request: %v", err)
			http.Error(w, fmt.Sprintf("Media upload failed: %v", err), http.StatusBadRequest)
			return
		}

		type response struct {
			URL string `json:"url"`
		}
		resp := response{
			URL: url,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
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

func mediaPath(contentType, username string) (string, error) {
	timestamp := dates.ThisFriday().Format("20060102")
	extension := "png"
	if contentType == "image/jpeg" {
		extension = "jpg"
	}

	return fmt.Sprintf("uploads/%s/%s/%s.%s", username, timestamp, randomKey(), extension), nil
}

func randomKey() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
