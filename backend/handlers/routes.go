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
	api.HandleFunc("/entries/following", s.requireAuthentication(s.entriesFollowingGet())).Methods(http.MethodGet)
	api.HandleFunc("/entries/{username}", s.entriesGet()).Methods(http.MethodGet)
	api.HandleFunc("/entries/{username}/project/{project}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/entries/{username}/project/{project}", s.projectGet()).Methods(http.MethodGet)
	api.HandleFunc("/entry/{date}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/entry/{date}", s.requireAuthentication(s.entryPut())).Methods(http.MethodPut)
	api.HandleFunc("/follow/{username}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/follow/{username}", s.requireAuthentication(s.followPut())).Methods(http.MethodPut)
	api.HandleFunc("/follow/{username}", s.requireAuthentication(s.followDelete())).Methods(http.MethodDelete)
	api.HandleFunc("/draft/{date}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/draft/{date}", s.requireAuthentication(s.draftGet())).Methods(http.MethodGet)
	api.HandleFunc("/draft/{date}", s.requireAuthentication(s.draftPut())).Methods(http.MethodPut)
	api.HandleFunc("/export", s.requireAuthentication(s.exportGet())).Methods(http.MethodGet)
	api.HandleFunc("/media", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/media", s.requireAuthentication(s.mediaPut())).Methods(http.MethodPut)
	api.HandleFunc("/pageViews", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/pageViews", s.pageViewsGet()).Methods(http.MethodGet)
	api.HandleFunc("/preferences", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/preferences", s.requireAuthentication(s.preferencesGet())).Methods(http.MethodGet)
	api.HandleFunc("/preferences", s.requireAuthentication(s.preferencesPut())).Methods(http.MethodPut)
	api.HandleFunc("/reactions/entry/{username}/{date}", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/reactions/entry/{username}/{date}", s.reactionsGet()).Methods(http.MethodGet)
	api.HandleFunc("/reactions/entry/{username}/{date}", s.requireAuthentication(s.reactionsPost())).Methods(http.MethodPost)
	api.HandleFunc("/reactions/entry/{username}/{date}", s.requireAuthentication(s.reactionsDelete())).Methods(http.MethodDelete)
	api.HandleFunc("/recentEntries", s.recentEntriesGet()).Methods(http.MethodGet)
	api.HandleFunc("/user/me", s.requireAuthentication(s.userMeGet())).Methods(http.MethodGet)
	api.HandleFunc("/user/avatar", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/user/avatar", s.requireAuthentication(s.userAvatarPut())).Methods(http.MethodPut)
	api.HandleFunc("/user/avatar", s.requireAuthentication(s.userAvatarDelete())).Methods(http.MethodDelete)
	api.HandleFunc("/user/{username}", s.userGet()).Methods(http.MethodGet)
	api.HandleFunc("/user/{username}/following", s.userFollowingGet()).Methods(http.MethodGet)
	api.HandleFunc("/user", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/user", s.requireAuthentication(s.userPost())).Methods(http.MethodPost)
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
	static.PathPrefix("/css/").HandlerFunc(s.serveStaticResource()).Methods(http.MethodGet)
	static.PathPrefix("/js/").HandlerFunc(s.serveStaticResource()).Methods(http.MethodGet)
	static.PathPrefix("/images/").HandlerFunc(s.serveStaticResource()).Methods(http.MethodGet)

	// Add all the root-level static resources.
	for _, f := range []string{
		"android-chrome-192x192.png",
		"android-chrome-256x256.png",
		"apple-touch-icon.png",
		"browserconfig.xml",
		"favicon-16x16.png",
		"favicon-32x32.png",
		"favicon.ico",
		"mstile-150x150.png",
		"robots.txt",
		"site.webmanifest",
	} {
		static.PathPrefix("/" + f).HandlerFunc(s.serveStaticResource()).Methods(http.MethodGet)
	}

	// Serve index.html, the base page HTML before Vue rendering happens, and
	// render certain page elements server-side.
	static.HandleFunc("/{username}/{date}", s.serveEntryOr404()).Methods(http.MethodGet)
	static.HandleFunc("/{username}/project/{project}", serveIndexPage).Methods(http.MethodGet)

	static.HandleFunc("/about", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/feed", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/recent", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/entry/edit/{date}", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/login", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/logout", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/preferences", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/privacy-policy", serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/profile/edit", serveIndexPage).Methods(http.MethodGet)

	static.HandleFunc("/{username}", s.serveUserProfileOr404()).Methods(http.MethodGet)
	static.HandleFunc("/", serveIndexPage).Methods(http.MethodGet)
	static.PathPrefix("/").HandlerFunc(serve404).Methods(http.MethodGet)
}
