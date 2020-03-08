package firestore

import (
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetPreferences retrieves the user's preferences for using the site.
func (c client) GetPreferences(username string) (types.Preferences, error) {
	doc := c.firestoreClient.Collection(preferencesRootKey).Doc(username)
	docsnap, err := doc.Get(c.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return types.Preferences{}, datastore.PreferencesNotFoundError{Username: username}
		}
		return types.Preferences{}, err
	}
	var prefs types.Preferences
	if err := docsnap.DataTo(&prefs); err != nil {
		return types.Preferences{}, err
	}
	return prefs, nil
}

// SetPreferences saves the user's preferences for using the site.
func (c client) SetPreferences(username string, prefs types.Preferences) error {
	log.Printf("saving preferences to datastore: %+v", prefs)
	_, err := c.firestoreClient.Collection(preferencesRootKey).Doc(username).Set(c.ctx, prefs)
	return err
}
