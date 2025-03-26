package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s defaultServer) reactionsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("invalid date: %s - %s", date, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entryAuthor, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		reactions, err := s.datastore.GetReactions(entryAuthor, date)
		if err != nil {
			log.Printf("failed to retrieve reactions: %s", err)
			http.Error(w, "Failed to retrieve reactions", http.StatusInternalServerError)
			return
		}

		// Check Accept header to determine response format
		acceptHeader := r.Header.Get("Accept")
		if strings.Contains(acceptHeader, "text/html") {
			// User wants HTML response
			var loggedInUsername types.Username
			// Try to get the username from context, but don't fail if not available
			if r.Context().Value(contextKeyUsername) != nil {
				loggedInUsername = mustGetUsernameFromContext(r.Context())
			}

			data := struct {
				EntryAuthor      string
				EntryDate        string
				LoggedInUsername string
				ReactionSymbols  []string
				Reactions        []types.Reaction
				SelectedReaction string
				CSRFToken        string
			}{
				EntryAuthor:      string(entryAuthor),
				EntryDate:        string(date),
				LoggedInUsername: string(loggedInUsername),
				ReactionSymbols:  []string{"👍", "🎉", "🙁"},
				Reactions:        reactions,
				CSRFToken:        csrf.Token(r),
			}

			// Find if user has a reaction already
			for _, reaction := range reactions {
				if reaction.Username == loggedInUsername {
					data.SelectedReaction = reaction.Symbol
					break
				}
			}

			// Render template
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := s.templates.ExecuteTemplate(w, "reactions", data); err != nil {
				log.Printf("error rendering reactions template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			// Default to JSON response (backward compatibility)
			respondOK(w, reactions)
		}
	}
}

func (s defaultServer) reactionsPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reactionSymbol, err := reactionSymbolFromRequest(r)
		if err != nil {
			log.Printf("invalid reactions request: %v", err)
			http.Error(w, "Invalid reactions request", http.StatusBadRequest)
			return
		}

		entryAuthor, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		entryDate, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("invalid date: %s", entryDate)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		loggedInUsername := mustGetUsernameFromContext(r.Context())
		reaction := types.Reaction{
			Username: loggedInUsername,
			Symbol:   reactionSymbol,
		}
		err = s.datastore.AddReaction(entryAuthor, entryDate, reaction)
		if err != nil {
			log.Printf("failed to add reaction: %s", err)
			http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
			return
		}

		// If client accepts HTML, return updated reactions HTML
		acceptHeader := r.Header.Get("Accept")
		if strings.Contains(acceptHeader, "text/html") {
			s.renderReactionsHTML(w, r, string(entryAuthor), string(entryDate), string(loggedInUsername))
		}
	}
}

func (s defaultServer) reactionsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entryAuthor, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		entryDate, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("invalid date: %s", entryDate)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		loggedInUsername := mustGetUsernameFromContext(r.Context())
		err = s.datastore.DeleteReaction(entryAuthor, entryDate, loggedInUsername)
		if err != nil {
			log.Printf("failed to delete reaction: %s", err)
			http.Error(w, "Failed to delete reaction", http.StatusInternalServerError)
			return
		}

		// If client accepts HTML, return updated reactions HTML
		acceptHeader := r.Header.Get("Accept")
		if strings.Contains(acceptHeader, "text/html") {
			s.renderReactionsHTML(w, r, string(entryAuthor), string(entryDate), string(loggedInUsername))
		}
	}
}

// Helper function to render reactions HTML after updates
func (s defaultServer) renderReactionsHTML(w http.ResponseWriter, r *http.Request, entryAuthor, entryDate, loggedInUsername string) {
	reactions, err := s.datastore.GetReactions(types.Username(entryAuthor), types.EntryDate(entryDate))
	if err != nil {
		log.Printf("failed to retrieve reactions: %s", err)
		http.Error(w, "Failed to retrieve reactions", http.StatusInternalServerError)
		return
	}

	data := struct {
		EntryAuthor      string
		EntryDate        string
		LoggedInUsername string
		ReactionSymbols  []string
		Reactions        []types.Reaction
		SelectedReaction string
		CSRFToken        string
	}{
		EntryAuthor:      entryAuthor,
		EntryDate:        entryDate,
		LoggedInUsername: loggedInUsername,
		ReactionSymbols:  []string{"👍", "🎉", "🙁"},
		Reactions:        reactions,
		CSRFToken:        csrf.Token(r),
	}

	// Find if user has a reaction already
	for _, reaction := range reactions {
		if reaction.Username == types.Username(loggedInUsername) {
			data.SelectedReaction = reaction.Symbol
			break
		}
	}

	// Render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "reactions", data); err != nil {
		log.Printf("error rendering reactions template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		return "", errors.New(`request is missing required field: "reactionSymbol"`)
	}

	reactionSymbol := *rr.ReactionSymbol
	if !isValidReaction(reactionSymbol) {
		return "", fmt.Errorf("invalid reaction choice: %s", reactionSymbol)
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
