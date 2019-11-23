package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/handlers/entry"
)

func (s *defaultServer) topicGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		topic, err := topicFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve topic from request path: %s", err)
			http.Error(w, "Invalid topic", http.StatusBadRequest)
			return
		}

		entries, err := s.datastore.GetEntries(username)
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			http.Error(w, fmt.Sprintf("Failed to retrieve entries for %s", username), http.StatusInternalServerError)
			return
		}

		type topicBody struct {
			Date     string `json:"date"`
			Markdown string `json:"markdown"`
		}

		topicBodies := []topicBody{}
		for _, e := range entries {
			body, err := entry.ReadTopic(e.Markdown, topic)
			if err != nil {
				continue
			}
			topicBodies = append(topicBodies, topicBody{
				Markdown: body,
				Date:     e.Date,
			})
		}

		if err := json.NewEncoder(w).Encode(topicBodies); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) topicOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
