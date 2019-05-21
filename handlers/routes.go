package handlers

import (
	"net/http"
)

func (s *defaultServer) routes() {
	fs := http.FileServer(http.Dir("./web/frontend/dist"))
	s.router.PathPrefix("/js").Handler(fs)
	s.router.PathPrefix("/css").Handler(fs)
	s.router.PathPrefix("/images").Handler(fs)

	s.router.HandleFunc("/api/entries/{username}", s.entriesHandler())
	s.router.HandleFunc("/api/entry/{username}/{date}", s.entryHandler())
	s.router.HandleFunc("/api/recentEntries", s.recentEntriesHandler())
	s.router.HandleFunc("/api/user/me", s.userMeHandler())
	s.router.HandleFunc("/api/submit", s.submitHandler())
	s.router.HandleFunc("/api/logout", s.logoutHandler())
	s.router.PathPrefix("/api").HandlerFunc(s.apiRootHandler())

	s.router.HandleFunc("/submit", s.submitPageHandler())
	s.router.HandleFunc("/me", s.meRedirectHandler())
	s.router.PathPrefix("/").HandlerFunc(s.indexHandler())
}
