package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func contentSecurityPolicy() string {
	directives := map[string][]string{
		"script-src": []string{
			"'self'",
			"https://www.google-analytics.com",
			"https://www.googletagmanager.com",
			// URLs for /login route (UserKit)
			"https://widget.userkit.io",
			"https://api.userkit.io",
			"https://www.google.com/recaptcha/",
			"https://www.gstatic.com/recaptcha/",
		},
		"style-src": []string{
			"'self'",
			// URLs for /login route (UserKit)
			"https://widget.userkit.io/css/",
			"https://fonts.googleapis.com",
			"https://fonts.gstatic.com",
		},
		"frame-src": []string{
			// URLs for /login route (UserKit)
			"https://www.google.com/recaptcha/",
		},
		"img-src": []string{
			"'self'",
			// For bootstrap navbar images
			"data:",
			// For Google Analytics
			"https://www.google-analytics.com",
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

func (s defaultServer) enableCsp(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", contentSecurityPolicy())
		h(w, r)
	}
}
