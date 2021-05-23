package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

// Many of these rules come from UserKit:
// https://docs.userkit.io/docs/content-security-policy
func contentSecurityPolicy() string {
	directives := map[string][]string{
		"script-src": {
			"'self'",
			"https://www.google-analytics.com",
			"https://www.googletagmanager.com",
			// URLs for /login route (UserKit)
			"https://widget.userkit.io",
			"https://api.userkit.io",
			"https://www.google.com/recaptcha/",
			"https://www.gstatic.com/recaptcha/",
			"https://apis.google.com",
		},
		"style-src": {
			"'self'",
			// URLs for /login route (UserKit)
			"https://widget.userkit.io",
			"https://fonts.googleapis.com",
			"https://fonts.gstatic.com",
			// Google auth requires this, and I can't figure out any way to avoid it.
			"'unsafe-inline'",
		},
		"frame-src": {
			// URLs for /login route (UserKit)
			"https://www.google.com/recaptcha/",
			"https://accounts.google.com",
		},
		"img-src": {
			"'self'",
			"https://storage.googleapis.com/whatgotdone-public/",
			// For bootstrap navbar images
			"data:",
			// For Google Analytics
			"https://www.google-analytics.com",
			// For Google Sign In
			"https://*.googleusercontent.com",
			// For UserKit
			"https://widget.userkit.io",
		},
	}
	directives["script-src"] = append(directives["script-src"], extraScriptSrcSources()...)
	directives["style-src"] = append(directives["style-src"], extraStyleSrcSources()...)
	policyParts := []string{}
	for directive, sources := range directives {
		policyParts = append(policyParts, fmt.Sprintf("%s %s", directive, strings.Join(sources, " ")))
	}
	return strings.Join(policyParts, "; ")
}

func enableCsp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", contentSecurityPolicy())
		next.ServeHTTP(w, r)
	})
}
