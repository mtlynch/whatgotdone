// +build !staging

package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func enableCsp(w *http.ResponseWriter) {
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
	frameSrc := strings.Join([]string{
		"'self'",
		// For sendinblue mailing list signup
		"https://sibforms.com",
	}, " ")
	(*w).Header().Set("Content-Security-Policy", fmt.Sprintf("default-src %s; img-src %s; frame-src %s", defaultSrc, imgSrc, frameSrc))
}
