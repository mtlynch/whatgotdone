package handlers

import (
	"net/http"
)

func (s *defaultServer) routes() {
	fs := http.FileServer(http.Dir("./web/frontend/dist"))
	s.router.PathPrefix("/js").Handler(fs)
	s.router.PathPrefix("/css").Handler(fs)
	s.router.PathPrefix("/images").Handler(fs)
	s.router.Path("/android-chrome-192x192.png").Handler(fs)
	s.router.Path("/android-chrome-512x512.png").Handler(fs)
	s.router.Path("/app.js").Handler(fs)
	s.router.Path("/apple-touch-icon.png").Handler(fs)
	s.router.Path("/browserconfig.xml").Handler(fs)
	s.router.Path("/favicon-16x16.png").Handler(fs)
	s.router.Path("/favicon-32x32.png").Handler(fs)
	s.router.Path("/favicon.ico").Handler(fs)
	s.router.Path("/mstile-150x150.png").Handler(fs)
	s.router.Path("/site.webmanifest").Handler(fs)

	s.router.HandleFunc("/api/entries/{username}", s.enableCors(s.entriesGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/draft/{date}", s.enableCors(s.draftOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/draft/{date}", s.enableCors(s.draftGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/draft/{date}", s.enableCors(s.draftPost())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.enableCors(s.reactionsOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.enableCors(s.reactionsGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.enableCors(s.reactionsPost())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/recentEntries", s.enableCors(s.recentEntriesGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user/me", s.enableCors(s.userMeGet())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/submit", s.enableCors(s.submitOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/submit", s.enableCors(s.submitPost())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/logout", s.enableCors(s.logoutOptions())).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/logout", s.enableCors(s.logoutPost())).Methods(http.MethodPost)
	s.router.PathPrefix("/api").HandlerFunc(s.enableCors(s.apiRootHandler()))

	s.router.HandleFunc("/submit", s.enableCsp(s.submitPageHandler()))
	s.router.PathPrefix("/").HandlerFunc(s.enableCsp(s.indexHandler()))
}
