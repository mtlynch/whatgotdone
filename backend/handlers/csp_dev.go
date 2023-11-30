//go:build dev || staging

// The dev environment needs extra CSP sources because when Vue generates
// frontend code in dev mode, it requires directives that would be unsafe in
// production.

package handlers

func extraScriptSrcSources() []string {
	return []string{"'unsafe-eval'", "'unsafe-inline'"}
}

func extraStyleSrcSources() []string {
	return []string{}
}

func extraImgSrcSources() []string {
	return []string{"https://storage.googleapis.com"}
}
