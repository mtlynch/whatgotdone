package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

var contextKeyCSPNonce = &contextKey{"csp-nonce"}

// Many of these rules come from UserKit:
// https://docs.userkit.io/docs/content-security-policy
func contentSecurityPolicy() string {
	directives := map[string][]string{
		"script-src": {
			"'self'",
			"https://plausible.io",
			// URLs for /login route (UserKit)
			"https://widget.userkit.io",
			"https://api.userkit.io",
			"https://www.google.com/recaptcha/",
			"https://www.gstatic.com/recaptcha/",
			"https://accounts.google.com",
		},
		"style-src": {
			"'self'",
			// URLs for /login route (UserKit)
			"https://accounts.google.com",
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
			// For bootstrap navbar images
			"data:",
			// For user-generated uploads
			"https://media.whatgotdone.com",
			// For Google Sign In
			"https://*.googleusercontent.com",
			// For UserKit
			"https://widget.userkit.io",
		},
	}
	directives["script-src"] = append(directives["script-src"], extraScriptSrcSources()...)
	directives["img-src"] = append(directives["img-src"], extraImgSrcSources()...)
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
		nonce := base64.StdEncoding.EncodeToString(random.Bytes(16))
		ctx := context.WithValue(r.Context(), contextKeyCSPNonce, nonce)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func cspNonce(ctx context.Context) string {
	key, ok := ctx.Value(contextKeyCSPNonce).(string)
	if !ok {
		panic("CSP nonce is missing from request context")
	}
	return key
}
