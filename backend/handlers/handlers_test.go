package handlers

import (
	"errors"
	"os"
	"path"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func init() {
	// The handler uses relative paths to find the template file. Switch to the
	// app's root directory so that the relative paths work.
	if err := os.Chdir("../../"); err != nil {
		panic(err)
	}
	frontendIndexPath := path.Join(frontendRootDir, frontendIndexFilename)

	// Ensure that the frontend/dist/index.html exists. The handler functions
	// need it, even if it's empty.
	if _, err := os.Stat(frontendIndexPath); os.IsNotExist(err) {
		// Ensure that the frontend/dist folder exists.
		if err = os.MkdirAll(frontendRootDir, os.ModePerm); err != nil {
			panic(err)
		}
		// Create frontend/dist/index.html.
		if _, err := os.Create(frontendIndexPath); err != nil {
			panic(err)
		}
	}
}

type mockAuthenticator struct {
	tokensToUsers map[string]types.Username
}

func (a mockAuthenticator) UserFromAuthToken(authToken string) (types.Username, error) {
	for k, v := range a.tokensToUsers {
		if k == authToken {
			return v, nil
		}
	}
	return "", errors.New("mock token not found")
}
