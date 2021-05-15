package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (s *defaultServer) userAvatarPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.gcsClient == nil {
			log.Printf("can't accept avatar upload because media uploads are disabled")
			http.Error(w, fmt.Sprintf("User profile photo uploading is disabled"), http.StatusBadRequest)
			return
		}

		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must be logged in to upload a profile photo", http.StatusForbidden)
			return
		}

		avatarFile, contentType, err := avatarFileFromRequest(w, r)
		if err != nil {
			log.Printf("failed to read media from request: %v", err)
			http.Error(w, fmt.Sprintf("Media upload failed: %v", err), http.StatusBadRequest)
			return
		}

		path := fmt.Sprintf("avatars/%s/%s-avatar.jpg", username, username)

		_, err = s.gcsClient.UploadFile(avatarFile, path, contentType)
		if err != nil {
			log.Printf("failed to upload avatar to static storage: %v", err)
			http.Error(w, fmt.Sprintf("Profile photo upload failed: %v", err), http.StatusBadRequest)
			return
		}
	}
}

const maxAvatarBytes = 8 * 1024 * 1024 // 8 MB

func avatarFileFromRequest(w http.ResponseWriter, r *http.Request) (io.Reader, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarBytes)
	r.ParseMultipartForm(32 << 20)
	file, metadata, err := r.FormFile("file")
	if err != nil {
		return nil, "", err
	}
	contentType := metadata.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return nil, "", errors.New("invalid avatar file format")
	}
	return file, contentType, nil
}
