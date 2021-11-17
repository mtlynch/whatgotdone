package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s defaultServer) reactionsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s - %s", date, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entryAuthor, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		reactions, err := s.datastore.GetReactions(entryAuthor, date)
		if err != nil {
			log.Printf("Failed to retrieve reactions: %s", err)
			http.Error(w, "Failed to retrieve reactions", http.StatusInternalServerError)
			return
		}

		reactionsFiltered := []types.Reaction{}
		for _, reaction := range reactions {
			if reaction.Symbol != "" {
				reactionsFiltered = append(reactionsFiltered, reaction)
			}
		}

		respondOK(w, reactionsFiltered)
	}
}

func (s defaultServer) reactionsPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to provide a reaction", http.StatusForbidden)
			return
		}

		reactionSymbol, err := reactionSymbolFromRequest(r)
		if err != nil {
			log.Printf("Invalid reactions request: %v", err)
			http.Error(w, "Invalid reactions request", http.StatusBadRequest)
			return
		}

		entryAuthor, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		entryDate, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", entryDate)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reaction := types.Reaction{
			Username:  username,
			Timestamp: time.Now().Format(time.RFC3339),
			Symbol:    reactionSymbol,
		}
		err = s.datastore.AddReaction(entryAuthor, entryDate, reaction)
		if err != nil {
			log.Printf("Failed to add reaction: %s", err)
			http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
			return
		}
	}
}

func (s defaultServer) reactionsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to delete a reaction", http.StatusForbidden)
			return
		}

		entryAuthor, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		entryDate, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", entryDate)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.datastore.DeleteReaction(entryAuthor, entryDate, username)
		if err != nil {
			log.Printf("Failed to delete reaction: %s", err)
			http.Error(w, "Failed to delete reaction", http.StatusInternalServerError)
			return
		}
	}
}

func reactionSymbolFromRequest(r *http.Request) (string, error) {
	type reactionRequest struct {
		ReactionSymbol *string `json:"reactionSymbol"`
	}
	var rr reactionRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rr)
	if err != nil {
		return "", err
	}

	if rr.ReactionSymbol == nil {
		return "", errors.New(`Request is missing required field: "reactionSymbol"`)
	}

	reactionSymbol := *rr.ReactionSymbol
	if !isValidReaction(reactionSymbol) {
		return "", fmt.Errorf("Invalid reaction choice: %s", reactionSymbol)
	}

	return reactionSymbol, nil
}

func isValidReaction(reaction string) bool {
	validReactionSymbols := [...]string{"", "👍", "🙁", "🎉"}
	for _, v := range validReactionSymbols {
		if reaction == v {
			return true
		}
	}
	return false
}
