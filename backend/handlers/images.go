package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func (s *defaultServer) imagesPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.gcsClient == nil {
			log.Printf("can't accept image upload because image uploads are disabled")
			http.Error(w, fmt.Sprintf("Image uploading is disabled"), http.StatusBadRequest)
			return
		}

		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must be logged in to upload an image", http.StatusForbidden)
			return
		}

		imageFile, contentType, err := imageFileFromRequest(w, r)
		if err != nil {
			log.Printf("failed to read image from request: %v", err)
			http.Error(w, fmt.Sprintf("Image upload failed: %v", err), http.StatusBadRequest)
			return
		}

		path, err := imagePath(contentType, username)
		if err != nil {
			log.Printf("failed to generate image path: %v", err)
			http.Error(w, fmt.Sprintf("Image upload failed: %v", err), http.StatusInternalServerError)
			return
		}

		url, err := s.gcsClient.UploadFile(imageFile, path, contentType)
		if err != nil {
			log.Printf("failed to read image from request: %v", err)
			http.Error(w, fmt.Sprintf("Image upload failed: %v", err), http.StatusBadRequest)
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

const maxImageSize = 20971520 // 20 MB

func imageFileFromRequest(w http.ResponseWriter, r *http.Request) (io.Reader, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize)
	r.ParseMultipartForm(32 << 20)
	file, metadata, err := r.FormFile("file")
	if err != nil {
		return nil, "", err
	}
	contentType := metadata.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return nil, "", errors.New("invalid image format")
	}
	return file, contentType, nil
}

func imagePath(contentType, username string) (string, error) {
	timestamp := thisFriday().Format("20060102")
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

// TODO: Remove duplicate version of this function.
func thisFriday() time.Time {
	t := time.Now()
	for t.Weekday() != time.Friday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}
