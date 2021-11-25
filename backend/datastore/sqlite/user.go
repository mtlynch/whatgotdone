package sqlite

import (
	"log"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetUserProfile returns profile information for the given user.
func (d db) GetUserProfile(username types.Username) (types.UserProfile, error) {
	stmt, err := d.ctx.Prepare(`
		SELECT
			about_markdown,
			email,
			twitter,
			mastodon
		FROM
			user_profiles
		WHERE
			username=?`)
	if err != nil {
		return types.UserProfile{}, err
	}
	defer stmt.Close()

	var (
		aboutMarkdown string
		email         string
		twitter       string
		mastodon      string
	)
	err = stmt.QueryRow(username).Scan(&aboutMarkdown, &email, &twitter, &mastodon)
	if err != nil {
		return types.UserProfile{}, err
	}

	return types.UserProfile{
		AboutMarkdown:   types.UserBio(aboutMarkdown),
		EmailAddress:    types.EmailAddress(email),
		TwitterHandle:   types.TwitterHandle(twitter),
		MastodonAddress: types.MastodonAddress(mastodon),
	}, nil
}

// SetUserProfile updates the given user's profile.
func (d db) SetUserProfile(username types.Username, profile types.UserProfile) error {
	log.Printf("saving user profile to datastore: %s -> %+v", username, profile)
	_, err := d.ctx.Exec(`
	INSERT OR REPLACE INTO user_profiles(
		username,
		about_markdown,
		email,
		twitter,
		mastodon)
	values(?,?,?,?,?)`, username, profile.AboutMarkdown, profile.EmailAddress, profile.TwitterHandle, profile.MastodonAddress)
	return err
}
