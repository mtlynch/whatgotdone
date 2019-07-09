package handlers

import (
	"net/http"
)

func (s *defaultServer) routes() {
	fs := http.FileServer(http.Dir("./web/frontend/dist"))
	s.router.PathPrefix("/js").Handler(fs)
	s.router.PathPrefix("/css").Handler(fs)
	s.router.PathPrefix("/images").Handler(fs)
	s.router.PathPrefix("/app.js").Handler(fs)

	s.router.HandleFunc("/api/entries/{username}", s.enableCors(s.entriesHandler()))
	s.router.HandleFunc("/api/draft/{date}", s.enableCors(s.draftHandler()))
	s.router.HandleFunc("/api/recentEntries", s.enableCors(s.recentEntriesHandler()))
	s.router.HandleFunc("/api/user/me", s.enableCors(s.userMeHandler()))
	s.router.HandleFunc("/api/submit", s.enableCors(s.submitHandler()))
	s.router.HandleFunc("/api/logout", s.enableCors(s.logoutHandler()))
	s.router.PathPrefix("/api").HandlerFunc(s.enableCors(s.apiRootHandler()))

	s.router.HandleFunc("/submit", s.enableCsp(s.submitPageHandler()))
	s.router.PathPrefix("/").HandlerFunc(s.enableCsp(s.indexHandler()))
}
