package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/image"
)

func (s *defaultServer) userAvatarGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("couldn't find username in avatar request path: %v", err)
			http.Error(w, "Couldn't find username in avatar request path", http.StatusBadRequest)
			return
		}
		avatarReader, err := s.datastore.GetAvatar(username)
		if err != nil {
			if _, ok := err.(datastore.ErrAvatarNotFound); ok {
				http.Redirect(w, r, "/images/no-avatar.jpg", http.StatusTemporaryRedirect)
				return
			}
			log.Printf("failed to read avatar: %v", err)
			http.Error(w, "Failed to read avatar", http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(w, avatarReader); err != nil {
			log.Printf("failed to write avatar in response: %v", err)
			http.Error(w, "Failed to copy avatar into response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *defaultServer) userAvatarPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.gcsClient == nil {
			log.Printf("can't accept avatar upload because GCS is disabled")
			http.Error(w, "User profile photo uploading is disabled", http.StatusBadRequest)
			return
		}

		avatarFile, contentType, err := avatarFileFromRequest(w, r)
		if err != nil {
			log.Printf("failed to read profile photo from request: %v", err)
			http.Error(w, fmt.Sprintf("Profile photo upload failed: %v", err), http.StatusBadRequest)
			return
		}

		if contentType != "image/jpeg" {
			http.Error(w, "Profile photo must be in JPEG format", http.StatusBadRequest)
			return
		}

		avatarRawImg, err := image.Decode(avatarFile, image.DecodeLimits{
			MinWidthPixels:  40,
			MaxWidthPixels:  50000,
			MinHeightPixels: 40,
			MaxHeightPixels: 50000,
		})
		if err != nil {
			log.Printf("failed to read profile photo from request: %v", err)
			http.Error(w, fmt.Sprintf("Profile photo upload failed: %v", err), http.StatusBadRequest)
			return
		}

		username := mustGetUsernameFromContext(r.Context())
		const avatarWidth = 300
		resizedAvatar := image.Resize(avatarRawImg, avatarWidth)
		var buf bytes.Buffer
		if err = image.Encode(resizedAvatar.Img, &buf); err != nil {
			log.Printf("failed to encode image to bytes: %v", err)
			http.Error(w, fmt.Sprintf("Profile photo upload failed: %v", err), http.StatusInternalServerError)
			return
		}
		if err := s.datastore.InsertAvatar(username, &buf, resizedAvatar.Width); err != nil {
			log.Printf("failed to insert avatar into database for user %s: %v", username, err)
			http.Error(w, fmt.Sprintf("Failed to save avatar for user %s", username), http.StatusInternalServerError)
			return
		}
	}
}

func (s *defaultServer) userAvatarDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mustGetUsernameFromContext(r.Context())

		if err := s.datastore.DeleteAvatar(username); err != nil {
			log.Printf("failed to delete avatar for user %s: %v", username, err)
			http.Error(w, fmt.Sprintf("failed to delete avatar for user %s", username), http.StatusInternalServerError)
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
