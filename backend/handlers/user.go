package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/handlers/validate"
	"github.com/mtlynch/whatgotdone/backend/types"
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
		if _, ok := err.(datastore.UserProfileNotFoundError); ok {
			http.Error(w, "No profile found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Failed to retrieve user profile data for %s: %s", username, err)
			http.Error(w, "Invalid username", http.StatusNotFound)
			return
		}

		type userResponse struct {
			AboutMarkdown string `json:"aboutMarkdown"`
			TwitterHandle string `json:"twitterHandle"`
			EmailAddress  string `json:"emailAddress"`
		}

		resp := userResponse{
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
			Username types.Username `json:"username"`
		}
		resp := userMeResponse{
			Username: username,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func profileFromRequest(r *http.Request) (types.UserProfile, error) {
	type profileUpdateRequest struct {
		AboutMarkdown string `json:"aboutMarkdown"`
		EmailAddress  string `json:"emailAddress"`
		TwitterHandle string `json:"twitterHandle"`
	}
	var pur profileUpdateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pur)
	if err != nil {
		return types.UserProfile{}, err
	}
	return types.UserProfile{
		AboutMarkdown: strings.TrimSpace(pur.AboutMarkdown),
		EmailAddress:  pur.EmailAddress,
		TwitterHandle: pur.TwitterHandle,
	}, nil
}

func (s defaultServer) userExists(username types.Username) (bool, error) {
	// BUG: Will only detect users who have published an entry. Ideally, we'd be
	// able to tell if the username exists in UserKit, but the UserKit API
	// currently does not offer lookup by username.
	users, err := s.datastore.Users()
	if err != nil {
		return false, err
	}
	for _, u := range users {
		if u == username {
			return true, nil
		}
	}
	return false, nil
}
