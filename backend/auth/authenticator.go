package auth

import (
	"log"
	"os"

	"github.com/mtlynch/whatgotdone/backend/types"
	userkit "github.com/workpail/userkit-go"
)

// Authenticator wraps a user authentication system.
type Authenticator interface {
	UserFromAuthToken(authToken string) (types.Username, error)
}

type (
	defaultAuthenticator struct {
		userKitClient userkit.Client
	}
)

// New creates a new Authenticator interface.
func New() Authenticator {
	sk := os.Getenv("USERKIT_SECRET")
	if sk == "" {
		log.Fatal("USERKIT_SECRET environment variable must be set")
	}
	return defaultAuthenticator{
		userKitClient: userkit.NewUserKit(sk),
	}
}

// UserFromAuthToken finds the user associated with the given auth token and
// returns that user's username.
func (a defaultAuthenticator) UserFromAuthToken(authToken string) (types.Username, error) {
	user, err := a.userKitClient.Users.GetUserBySession(authToken)
	if err != nil {
		log.Printf("Failed to authenticate user's session token with UserKit: %v", err)
		return "", err
	}
	return types.Username(user.Username), nil
}
