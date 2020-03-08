package handlers

import (
	"net/http"
)

// A no-op function that tells the router to accept the OPTIONS method for a
// particular route.
func allowOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s *defaultServer) routes() {
	s.router.Use(s.enableCors)
	s.router.Use(s.enableCsrf)

	// Handle routes that require backend logic.
	s.router.HandleFunc("/api/entries/following", s.entriesFollowingGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/entries/{username}", s.entriesGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/entries/{username}/project/{project}", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/entries/{username}/project/{project}", s.projectGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/entry/{date}", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/entry/{date}", s.entryPost()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/follow/{username}", followOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/follow/{username}", s.followPut()).Methods(http.MethodPut)
	s.router.HandleFunc("/api/follow/{username}", s.followDelete()).Methods(http.MethodDelete)
	s.router.HandleFunc("/api/draft/{date}", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/draft/{date}", s.draftGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/draft/{date}", s.draftPost()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/pageViews", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/pageViews", s.pageViewsGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.reactionsGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/reactions/entry/{username}/{date}", s.reactionsPost()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/recentEntries", s.recentEntriesGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user/me", s.userMeGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user/{username}", s.userGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user/{username}/following", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/user/{username}/following", s.userFollowingGet()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/user", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/user", s.userPost()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/logout", allowOptions()).Methods(http.MethodOptions)
	s.router.HandleFunc("/api/logout", s.logoutPost()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/tasks/refreshGoogleAnalytics", s.refreshGoogleAnalytics()).Methods(http.MethodGet)

	// Catchall for when no API route matches.
	s.router.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	})

	s.router.HandleFunc("/sitemap.xml", s.sitemapGet()).Methods(http.MethodGet)

	// Serve index.html, the base page HTML before Vue rendering happens, and
	// render certain page elements server-side.
	s.router.PathPrefix("/{username}/{date}").HandlerFunc(s.enableCsp(s.serveStaticResource())).Methods(http.MethodGet)
	s.router.PathPrefix("/").HandlerFunc(s.enableCsp(s.serveStaticResource())).Methods(http.MethodGet)
}
