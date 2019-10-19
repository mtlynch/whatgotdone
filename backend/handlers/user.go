package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s defaultServer) userGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		p, err := s.datastore.GetUserProfile(username)
		if err != nil {
			log.Printf("Failed to retrieve user profile data for %s: %s", username, err)
			http.Error(w, "Invalid username", http.StatusNotFound)
			return
		}

		type userResponse struct {
			Username      string `json:"username"`
			AboutMarkdown string `json:"aboutMarkdown"`
			TwitterHandle string `json:"twitterHandle"`
			EmailAddress  string `json:"emailAddress"`
		}

		resp := userResponse{
			Username:      username,
			AboutMarkdown: p.AboutMarkdown,
			TwitterHandle: p.TwitterHandle,
			EmailAddress:  p.EmailAddress,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func (s defaultServer) userMeGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must be logged in to retrieve information about your account", http.StatusForbidden)
			return
		}

		type userMeResponse struct {
			Username string `json:"username"`
		}

		resp := userMeResponse{
			Username: username,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
