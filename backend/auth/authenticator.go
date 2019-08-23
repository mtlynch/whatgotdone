package auth

import (
	"log"
	"os"

	userkit "github.com/workpail/userkit-go"
)

type Authenticator interface {
	UserFromAuthToken(authToken string) (string, error)
}

type (
	defaultAuthenticator struct {
		userKitClient userkit.Client
	}
)

func New() Authenticator {
	sk := os.Getenv("USERKIT_SECRET")
	if sk == "" {
		log.Fatal("USERKIT_SECRET environment variable must be set")
	}
	return defaultAuthenticator{
		userKitClient: userkit.NewUserKit(sk),
	}
}

func (a defaultAuthenticator) UserFromAuthToken(authToken string) (string, error) {
	user, err := a.userKitClient.Users.GetUserBySession(authToken)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}
