package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (s defaultServer) userAvatarFullGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}
		// TODO: Replace placeholder with this
		url := fmt.Sprintf("https://storage.googleapis.com/%s/avatars/%s/full-300px.jpg", "whatgotdone-staging", username)

		if username == "testjoe" {
			url = "https://placekitten.com/300/300"
		}

		http.Redirect(
			w,
			r,
			url,
			http.StatusTemporaryRedirect)
	}
}

func (s defaultServer) userAvatarThumbnailGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}
		// TODO: Replace placeholder with this
		url := fmt.Sprintf("https://storage.googleapis.com/%s/avatars/%s/thumb-30px.jpg", "whatgotdone-staging", username)

		if username == "testjoe" {
			url = "https://placekitten.com/300/300"
		}

		http.Redirect(
			w,
			r,
			url,
			http.StatusTemporaryRedirect)
	}
}
