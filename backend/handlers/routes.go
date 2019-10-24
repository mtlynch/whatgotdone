package handlers

import (
	"net/http"
)

func (s *defaultServer) routes() {
	s.router.Use(s.enableCors)

	// Handle routes that require backend logic.
	s.router.HandleFunc("/api/entries/{username}", s.entriesGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/draft/{date}", s.draftOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/draft/{date}", s.draftGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/draft/{date}", s.enableCsrf(s.draftPost())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.reactionsOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.reactionsGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.enableCsrf(s.reactionsPost())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/recentEntries", s.recentEntriesGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user/me", s.userMeGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user/{username}", s.userGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user", s.userOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/user", s.userPost()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/submit", s.submitOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/submit", s.enableCsrf(s.submitPost())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/logout", s.logoutOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/logout", s.enableCsrf(s.logoutPost())).Methods(http.MethodPost)
	// Catchall for when no API route matches.
	s.router.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	})

	// Serve index.html, the base page HTML before Vue rendering happens, and
	// render certain page elements server-side.
	s.router.PathPrefix("/{username}/{date}").HandlerFunc(s.enableCsrf(s.enableCsp(s.serveStaticResource()))).Methods(http.MethodGet)
	s.router.PathPrefix("/").HandlerFunc(s.enableCsrf(s.enableCsp(s.serveStaticResource()))).Methods(http.MethodGet)
}
