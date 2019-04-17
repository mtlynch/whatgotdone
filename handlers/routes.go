package handlers

import (
	"net/http"
)

func (s *defaultServer) routes() {
	fs := http.FileServer(http.Dir("./web/frontend/dist"))
	s.router.PathPrefix("/js").Handler(fs)
	s.router.PathPrefix("/css").Handler(fs)

	s.router.HandleFunc("/entries", s.entriesHandler)
	s.router.HandleFunc("/api/submit", s.submitHandler)
	s.router.PathPrefix("/").HandlerFunc(s.indexHandler)
}
