package auth

import (
	userkit "github.com/workpail/userkit-go"

	"github.com/mtlynch/whatgotdone/datastore"
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
	ks := datastore.NewUserKitKeyStore()
	sk, err := ks.SecretKey()
	if err != nil {
		panic(err)
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
