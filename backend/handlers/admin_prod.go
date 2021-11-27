// +build !dev
// +build !integration

package handlers

import "github.com/mtlynch/whatgotdone/backend/types"

// This is a temporary hack while we migrate from AppEngine to fly.io.
func isAdminUser(username types.Username) bool {
	return username == "michael"
}
