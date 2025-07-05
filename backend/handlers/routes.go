package handlers

import (
	"net/http"
)

func (s *defaultServer) routes() {
	s.router.Use(s.populateAuthenticationContext)

	authenticatedApis := s.router.PathPrefix("/api").Subrouter()
	authenticatedApis.Use(s.requireAuthenticationForApi)
	authenticatedApis.Use(s.enableCors)
	authenticatedApis.Use(s.enableCsrf)
	authenticatedApis.HandleFunc("/entries/following", s.entriesFollowingGet()).Methods(http.MethodGet)
	authenticatedApis.HandleFunc("/entry/{date}", s.entryPut()).Methods(http.MethodPut)
	authenticatedApis.HandleFunc("/entry/{date}", s.entryDelete()).Methods(http.MethodDelete)
	authenticatedApis.HandleFunc("/follow/{username}", s.followPut()).Methods(http.MethodPut)
	authenticatedApis.HandleFunc("/follow/{username}", s.followDelete()).Methods(http.MethodDelete)
	authenticatedApis.HandleFunc("/draft/{date}", s.draftGet()).Methods(http.MethodGet)
	authenticatedApis.HandleFunc("/draft/{date}", s.draftPut()).Methods(http.MethodPut)
	authenticatedApis.HandleFunc("/draft/{date}", s.draftDelete()).Methods(http.MethodDelete)
	authenticatedApis.HandleFunc("/export", s.exportGet()).Methods(http.MethodGet)
	authenticatedApis.HandleFunc("/export/markdown", s.exportMarkdownGet()).Methods(http.MethodGet)
	authenticatedApis.HandleFunc("/media", s.mediaPut()).Methods(http.MethodPut)
	authenticatedApis.HandleFunc("/preferences", s.preferencesGet()).Methods(http.MethodGet)
	authenticatedApis.HandleFunc("/preferences", s.preferencesPut()).Methods(http.MethodPut)
	authenticatedApis.HandleFunc("/reactions/entry/{username}/{date}", s.reactionsPost()).Methods(http.MethodPost)
	authenticatedApis.HandleFunc("/reactions/entry/{username}/{date}", s.reactionsDelete()).Methods(http.MethodDelete)
	authenticatedApis.HandleFunc("/user/me", s.userMeGet()).Methods(http.MethodGet)
	authenticatedApis.HandleFunc("/user/avatar", s.userAvatarPut()).Methods(http.MethodPut)
	authenticatedApis.HandleFunc("/user/avatar", s.userAvatarDelete()).Methods(http.MethodDelete)
	authenticatedApis.HandleFunc("/user", s.userPost()).Methods(http.MethodPost)

	apis := s.router.PathPrefix("/api").Subrouter()
	apis.Use(s.enableCors)
	apis.Use(s.enableCsrf)
	apis.HandleFunc("/entries/{username}", s.entriesGet()).Methods(http.MethodGet)
	apis.HandleFunc("/entries/{username}/project/{project}", s.projectGet()).Methods(http.MethodGet)
	apis.HandleFunc("/reactions/entry/{username}/{date}", s.reactionsGet()).Methods(http.MethodGet)
	apis.HandleFunc("/recentEntries", s.recentEntriesGet()).Methods(http.MethodGet)
	apis.HandleFunc("/user/{username}", s.userGet()).Methods(http.MethodGet)
	apis.HandleFunc("/user/{username}/following", s.userFollowingGet()).Methods(http.MethodGet)
	apis.HandleFunc("/logout", s.logoutPost()).Methods(http.MethodPost)

	s.addDevRoutes(apis)

	// Catchall for when no API route matches.
	apis.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	})

	authenticatedView := s.router.PathPrefix("/").Subrouter()
	authenticatedView.Use(upgradeToHttps)
	authenticatedView.Use(enableCsp)
	authenticatedView.Use(s.enableCsrf)
	authenticatedView.Use(s.requireAuthenticationForView)
	authenticatedView.HandleFunc("/entry/edit/{date}", s.serveIndexPage).Methods(http.MethodGet)
	authenticatedView.HandleFunc("/preferences", s.serveIndexPage).Methods(http.MethodGet)
	authenticatedView.HandleFunc("/feed", s.serveIndexPage).Methods(http.MethodGet)
	authenticatedView.HandleFunc("/profile/edit", s.serveIndexPage).Methods(http.MethodGet)

	static := s.router.PathPrefix("/").Subrouter()
	static.Use(upgradeToHttps)
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
	static.HandleFunc("/{username}/avatar", s.userAvatarGet()).Methods(http.MethodGet)

	// Serve index.html, the base page HTML before Vue rendering happens, and
	// render certain page elements server-side.
	static.HandleFunc("/{username}/{date}", s.serveEntryOr404()).Methods(http.MethodGet)
	static.HandleFunc("/{username}/project/{project}", s.serveIndexPage).Methods(http.MethodGet)

	static.HandleFunc("/about", s.serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/recent", s.serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/login", s.serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/logout", s.serveIndexPage).Methods(http.MethodGet)
	static.HandleFunc("/privacy-policy", s.serveIndexPage).Methods(http.MethodGet)

	static.HandleFunc("/{username}", s.serveUserProfileOr404()).Methods(http.MethodGet)
	static.HandleFunc("/", s.serveIndexPage).Methods(http.MethodGet)
	static.PathPrefix("/").HandlerFunc(s.serve404).Methods(http.MethodGet)

}
