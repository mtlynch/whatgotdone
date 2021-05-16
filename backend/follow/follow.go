package entries

import (
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// Store reads datastore content related to user following relationships.
type Store interface {
	// Users returns all the users who have published entries.
	Users() ([]types.Username, error)
	// InsertFollow adds a following relationship to the datastore.
	InsertFollow(leader, follower types.Username) error
}

// Manager stores information related to journal entries.
type Manager interface {
	Follow(follower, leader string) error
}

type defaultManager struct {
	store Store
}

func (m defaultManager) Follow(follower, leader types.Username) error {
	if ok, err := m.userExists(leader); !ok {
		log.Printf("user %s tried to follow non-existent user: %s", follower, leader)
		http.Error(w, "Invalid username", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("failed to look up whether user exists: %s - %s", leader, err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	err = s.datastore.InsertFollow(leader, follower)
	if _, ok := err.(datastore.FollowAlreadyExistsError); ok {
		log.Printf("tried to re-follow when %s -> %s when follow already existed", follower, leader)
		http.Error(w, "You're already following this user", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("failed to add follower: %s->%s - %s", follower, leader, err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}
}

func (m defaultManager) userExists(username types.Username) (bool, error) {
	// BUG: Will only detect users who have published an entry. Ideally, we'd be
	// able to tell if the username exists in UserKit, but the UserKit API
	// currently does not offer lookup by username.
	users, err := m.store.Users()
	if err != nil {
		return false, err
	}
	for _, u := range users {
		if u == username {
			return true, nil
		}
	}
	return false, nil
}

func New(store Store) Manager {
	return defaultManager{
		store: store,
	}
}
