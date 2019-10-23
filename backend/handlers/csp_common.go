package handlers

import (
	"fmt"
	"strings"
)

func contentSecurityPolicy() string {
	scriptSrcElem := strings.Join([]string{
		"'self'",
		"https://www.google-analytics.com",
		"https://www.googletagmanager.com",
		// URLs for /login route (UserKit)
		"https://widget.userkit.io",
		"https://api.userkit.io",
		"https://www.google.com/recaptcha/",
		"https://www.gstatic.com/recaptcha/",
	}, " ")
	styleSrcElem := strings.Join([]string{
		"'self'",
		// URLs for /login route (UserKit)
		"https://widget.userkit.io/css/",
		"https://fonts.googleapis.com",
		"https://fonts.gstatic.com",
	}, " ")
	imgSrc := strings.Join([]string{
		"'self'",
		// For bootstrap navbar images
		"data:",
		// For Google Analytics
		"https://www.google-analytics.com",
	},
		" ")
	return fmt.Sprintf("script-src-elem %s; style-src-elem %s; img-src %s", scriptSrcElem, styleSrcElem, imgSrc)
}
