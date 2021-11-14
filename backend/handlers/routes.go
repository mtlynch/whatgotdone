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
	api := s.router.PathPrefix("/api").Subrouter()
	api.Use(s.enableCors)
	api.Use(s.enableCsrf)
	api.HandleFunc("/entries/following", s.entriesFollowingGet()).Methods(http.MethodGet)
	api.HandleFunc("/entries/{username}", s.entriesGet()).Methods(http.MethodGet)
	api.HandleFunc("/entries/{username}/project/{project}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/entries/{username}/project/{project}", s.projectGet()).Methods(http.MethodGet)
	api.HandleFunc("/entry/{date}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/entry/{date}", s.entryPut()).Methods(http.MethodPut)
	api.HandleFunc("/follow/{username}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/follow/{username}", s.followPut()).Methods(http.MethodPut)
	api.HandleFunc("/follow/{username}", s.followDelete()).Methods(http.MethodDelete)
	api.HandleFunc("/draft/{date}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/draft/{date}", s.draftGet()).Methods(http.MethodGet)
	api.HandleFunc("/draft/{date}", s.draftPut()).Methods(http.MethodPut)
	api.HandleFunc("/media", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/media", s.mediaPut()).Methods(http.MethodPut)
	api.HandleFunc("/pageViews", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/pageViews", s.pageViewsGet()).Methods(http.MethodGet)
	api.HandleFunc("/preferences", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/preferences", s.preferencesGet()).Methods(http.MethodGet)
	api.HandleFunc("/preferences", s.preferencesPost()).Methods(http.MethodPost)
	api.HandleFunc("/reactions/entry/{username}/{date}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/reactions/entry/{username}/{date}", s.reactionsGet()).Methods(http.MethodGet)
	api.HandleFunc("/reactions/entry/{username}/{date}", s.reactionsPost()).Methods(http.MethodPost)
	api.HandleFunc("/recentEntries", s.recentEntriesGet()).Methods(http.MethodGet)
	api.HandleFunc("/user/me", s.userMeGet()).Methods(http.MethodGet)
	api.HandleFunc("/user/avatar", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/user/avatar", s.userAvatarPut()).Methods(http.MethodPut)
	api.HandleFunc("/user/avatar", s.userAvatarDelete()).Methods(http.MethodDelete)
	api.HandleFunc("/user/{username}", s.userGet()).Methods(http.MethodGet)
	api.HandleFunc("/user/{username}/following", s.userFollowingGet()).Methods(http.MethodGet)
	api.HandleFunc("/user", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/user", s.userPost()).Methods(http.MethodPost)
	api.HandleFunc("/logout", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/logout", s.logoutPost()).Methods(http.MethodPost)
	api.HandleFunc("/tasks/refreshGoogleAnalytics", s.refreshGoogleAnalytics()).Methods(http.MethodGet)

	// Catchall for when no API route matches.
	api.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	})

	static := s.router.PathPrefix("/").Subrouter()
	static.Use(enableCsp)
	static.Use(s.enableCsrf)
	static.HandleFunc("/sitemap.xml", s.sitemapGet()).Methods(http.MethodGet)

	// Serve index.html, the base page HTML before Vue rendering happens, and
	// render certain page elements server-side.
	static.HandleFunc("/{username}/{date}", s.serveStaticResource()).Methods(http.MethodGet)
	static.PathPrefix("/").HandlerFunc(s.serveStaticResource()).Methods(http.MethodGet)
}
