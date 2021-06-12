package parse

// If something goes wrong in a JavaScript-based client, it will send the
// literal string "undefined" as the when a variable is undefined.
func isUndefinedFromJavascript(s string) bool {
	return s == "undefined"
}
