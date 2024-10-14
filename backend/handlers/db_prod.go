//go:build !dev && !staging

package handlers

import (
	"github.com/gorilla/mux"
)

func (s *defaultServer) addDevRoutes(router *mux.Router) {
	// no-op
}
