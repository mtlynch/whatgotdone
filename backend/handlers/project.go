package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/handlers/entry"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s *defaultServer) projectGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		project, err := projectFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve project from request path: %s", err)
			http.Error(w, "Invalid project", http.StatusBadRequest)
			return
		}

		entries, err := s.datastore.GetEntries(username)
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			http.Error(w, fmt.Sprintf("Failed to retrieve entries for %s", username), http.StatusInternalServerError)
			return
		}

		type projectBody struct {
			Date     types.EntryDate `json:"date"`
			Markdown string          `json:"markdown"`
		}

		projectBodies := []projectBody{}
		for _, e := range entries {
			body, err := entry.ReadProject(e.Markdown, project)
			if _, ok := err.(entry.ProjectNotFoundError); ok {
				continue
			} else if err != nil {
				log.Printf("Failed to retrieve project from entry: %s", err)
				continue
			} else if body == "" {
				continue
			}
			projectBodies = append(projectBodies, projectBody{
				Markdown: body,
				Date:     e.Date,
			})
		}

		if err := json.NewEncoder(w).Encode(projectBodies); err != nil {
			panic(err)
		}
	}
}
