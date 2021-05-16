package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s defaultServer) followPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		follower, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to follow a new user", http.StatusForbidden)
			return
		}

		leader, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		if leader == follower {
			http.Error(w, "You can't follow yourself", http.StatusBadRequest)
			return
		}

		if ok, err := s.userExists(leader); !ok {
			log.Printf("user %s tried to follow non-existent user: %s", follower, leader)
			http.Error(w, "Invalid username", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("failed to look up whether user exists: %s - %s", leader, err)
			http.Error(w, "Failed to follow user", http.StatusInternalServerError)
			return
		}

		err = s.datastore.InsertFollow(leader, follower)
		if _, ok := err.(datastore.FollowAlreadyExistsError); ok {
			log.Printf("tried to re-follow when %s -> %s when follow already existed", follower, leader)
			http.Error(w, "You're already following this user", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to add follower: %s->%s - %s", follower, leader, err)
			http.Error(w, "Failed to follow user", http.StatusInternalServerError)
			return
		}
	}
}

func (s defaultServer) followDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		follower, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to unfollow a user", http.StatusForbidden)
			return
		}

		leader, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		err = s.datastore.DeleteFollow(leader, follower)
		if err != nil {
			log.Printf("failed to delete follower: %s->%s - %s", follower, leader, err)
			http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
			return
		}
	}
}

func (s defaultServer) userFollowingGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		follower, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		leaders, err := s.datastore.Following(follower)
		if err != nil {
			log.Printf("failed to find following list for %s: %v", follower, err)
			http.Error(w, fmt.Sprintf("Failed to retrieve users %s is following", follower), http.StatusInternalServerError)
			return
		}

		type response struct {
			Following []types.Username `json:"following"`
		}
		resp := response{
			Following: leaders,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
