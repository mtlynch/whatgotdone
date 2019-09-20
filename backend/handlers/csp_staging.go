// +build staging

package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func (s defaultServer) enableCsp(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defaultSrc := strings.Join([]string{
			"'self'",
			// URLs for /login route
			"https://widget.userkit.io",
			"https://api.userkit.io",
			"https://www.google.com/recaptcha/api.js",
			"https://www.gstatic.com/recaptcha/api2/",
			"https://fonts.googleapis.com",
			"https://fonts.gstatic.com",
			"https://www.google-analytics.com",
			"https://www.googletagmanager.com",
		}, " ")
		imgSrc := strings.Join([]string{
			"'self'",
			// For bootstrap navbar images
			"data:",
			// For Google Analytics
			"https://www.google-analytics.com",
		}, " ")
		w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src %s; img-src %s; unsafe-eval", defaultSrc, imgSrc, frameSrc))

		h(w, r)
	}
}
