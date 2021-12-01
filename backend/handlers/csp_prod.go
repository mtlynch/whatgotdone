// +build !dev

package handlers

func extraScriptSrcSources() []string {
	return []string{"'unsafe-eval'"}
}

func extraStyleSrcSources() []string {
	return []string{}
}
