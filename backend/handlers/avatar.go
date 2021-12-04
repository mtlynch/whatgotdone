package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/image"
)

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

		username := usernameFromContext(r.Context())
		const avatarThumbnailWidth = 40
		const avatarLargeWidth = 300
		for _, resizedAvatar := range image.Resize(avatarRawImg, []int{avatarLargeWidth, avatarThumbnailWidth}) {
			var buf bytes.Buffer
			err = image.Encode(resizedAvatar.Img, &buf)
			if err != nil {
				log.Printf("failed to encode image to bytes: %v", err)
				http.Error(w, fmt.Sprintf("Profile photo upload failed: %v", err), http.StatusInternalServerError)
				return
			}
			path := fmt.Sprintf("avatars/%s/%s-avatar-%dpx.jpg", username, username, resizedAvatar.Width)
			_, err = s.gcsClient.UploadFile(&buf, path, "image/jpeg", "no-cache")
			if err != nil {
				log.Printf("failed to upload avatar to static storage: %v", err)
				http.Error(w, fmt.Sprintf("Profile photo upload failed: %v", err), http.StatusInternalServerError)
				return
			}
		}

	}
}

func (s *defaultServer) userAvatarDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.gcsClient == nil {
			log.Printf("can't accept avatar delete because GCS is disabled")
			http.Error(w, "User profile photo deleting is disabled", http.StatusBadRequest)
			return
		}

		username := usernameFromContext(r.Context())

		path := fmt.Sprintf("avatars/%s/", username)
		if err := s.gcsClient.DeletePath(path); err != nil {
			log.Printf("failed to delete avatars from static storage: %v", err)
			http.Error(w, fmt.Sprintf("Profile photo deletion failed: %v", err), http.StatusInternalServerError)
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
