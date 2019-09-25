package handlers

import (
	"net/http"
)

// addAPI will add the api routes that require backend logic.
func (s *defaultServer) addAPI() {
	// API:entries
	s.router.HandleFunc("/api/entries/{username}", s.enableCors(s.entriesGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/recentEntries", s.enableCors(s.recentEntriesGet())).Methods(http.MethodGet)
	// API:draft
	s.router.HandleFunc("/api/draft/{date}", s.enableCors(s.draftOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/draft/{date}", s.enableCors(s.draftGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/draft/{date}", s.enableCsrf(s.enableCors(s.draftPost()))).Methods(http.MethodPost)
	// API:reactions
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.enableCors(s.reactionsOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.enableCors(s.reactionsGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.enableCsrf(s.enableCors(s.reactionsPost()))).Methods(http.MethodPost)
	// API:submit
	s.router.HandleFunc("/api/submit", s.enableCors(s.submitOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/submit", s.enableCsrf(s.enableCors(s.submitPost()))).Methods(http.MethodPost)
	// API:user
	s.router.HandleFunc("/api/user/me", s.enableCors(s.userMeGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/logout", s.enableCors(s.logoutOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/logout", s.enableCsrf(s.enableCors(s.logoutPost()))).Methods(http.MethodPost)
	// Catchall for when no API route matches.
	s.router.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	})
}

// addStatic will add the static management
func (s *defaultServer) addStatic() {
	// TODO make this location configurable in the future with flags or env.
	fs := http.FileServer(staticSystem{http.Dir(http.Dir("./frontend/dist"))})

	s.router.
		PathPrefix("/").
		HandlerFunc(s.enableCsrf(s.enableCsp(fs.ServeHTTP))).
		Methods(http.MethodGet)
}

func (s *defaultServer) routes() {
	// Add static file management
	s.addStatic()
	// Add the API
	s.addAPI()
	// render certain page elements server-side.
	s.router.PathPrefix("/{username}/{date}").HandlerFunc(s.enableCsrf(s.enableCsp(s.indexHandler()))).Methods(http.MethodGet)
}
