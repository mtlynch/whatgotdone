package handlers

import "net/http"

func requireTaskWorkerAuthentication(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Authenticate that this request came from an authenticated What Got
		// Done task worker.
		// We're skipping this right now because the worst thing an unauthenticated
		// user could do here is force us to query Google Analytics too frequently.
		fn(w, r)
	}
}
