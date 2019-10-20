package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/handlers/validate"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s defaultServer) userOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s defaultServer) userGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		p, err := s.datastore.GetUserProfile(username)

		if _, ok := err.(datastore.UserProfileNotFoundError); ok {
			http.Error(w, "No profile found", http.StatusNotFound)
			return
		} else if err != nil {
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

func (s defaultServer) userPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to update your profile", http.StatusForbidden)
			return
		}

		userProfile, err := profileFromRequest(r)
		if err != nil {
			log.Printf("Invalid profile update request: %v", err)
			http.Error(w, "Invalid profile update request", http.StatusBadRequest)
			return
		}

		if !validate.UserBio(userProfile.AboutMarkdown) {
			http.Error(w, "Invalid user bio", http.StatusBadRequest)
			return
		}

		if userProfile.EmailAddress != "" && !validate.EmailAddress(userProfile.EmailAddress) {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		if userProfile.TwitterHandle != "" && !validate.TwitterHandle(userProfile.TwitterHandle) {
			http.Error(w, "Invalid twitter handle", http.StatusBadRequest)
			return
		}

		err = s.datastore.SetUserProfile(username, userProfile)
		if err != nil {
			log.Printf("Failed to update user profile: %s", err)
			http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
			return
		}

		type profileUpdateResponse struct {
			Ok bool `json:"ok"`
		}
		resp := profileUpdateResponse{
			Ok: true,
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

type profileUpdateRequest struct {
	AboutMarkdown string `json:"aboutMarkdown"`
	EmailAddress  string `json:"emailAddress"`
	TwitterHandle string `json:"twitterHandle"`
}

func profileFromRequest(r *http.Request) (types.UserProfile, error) {
	var pur profileUpdateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pur)
	if err != nil {
		return types.UserProfile{}, err
	}
	return types.UserProfile{
		AboutMarkdown: pur.AboutMarkdown,
		EmailAddress:  pur.EmailAddress,
		TwitterHandle: pur.TwitterHandle,
	}, nil
}

func isValidProfile(profileUpdateRequest string) bool {
	return false
}
