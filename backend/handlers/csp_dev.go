//go:build dev

// The dev environment needs extra CSP sources because when Vue generates
// frontend code in dev mode, it requires directives that would be unsafe in
// production.

package handlers

func extraScriptSrcSources() []string {
	return []string{"'unsafe-eval'"}
}

func extraStyleSrcSources() []string {
	return []string{}
}
