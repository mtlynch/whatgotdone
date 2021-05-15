package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s defaultServer) userAvatarGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		type response struct {
			Url string `json:"url"`
		}
		resp := response{
			Url: "https://whatgotdone.com/" + username,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
