package sqlite

import (
	"database/sql"
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetUserProfile returns profile information for the given user.
func (d DB) GetUserProfile(username types.Username) (types.UserProfile, error) {
	var (
		aboutMarkdown string
		email         string
		twitter       string
		mastodon      string
	)
	err := d.ctx.QueryRow(`
		SELECT
				about_markdown,
				email,
				twitter,
				mastodon
		FROM
				user_profiles
		WHERE
				username = :username`,
		sql.Named("username", username)).Scan(&aboutMarkdown, &email, &twitter, &mastodon)

	if err == sql.ErrNoRows {
		return types.UserProfile{}, datastore.UserProfileNotFoundError{
			Username: username,
		}
	} else if err != nil {
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
func (d DB) SetUserProfile(username types.Username, profile types.UserProfile) error {
	log.Printf("saving user profile to datastore: %s -> %+v", username, profile)
	_, err := d.ctx.Exec(`
		INSERT OR REPLACE INTO user_profiles(
				username,
				about_markdown,
				email,
				twitter,
				mastodon)
		values(:username, :about_markdown, :email, :twitter, :mastodon)`,
		sql.Named("username", username),
		sql.Named("about_markdown", profile.AboutMarkdown),
		sql.Named("email", profile.EmailAddress),
		sql.Named("twitter", profile.TwitterHandle),
		sql.Named("mastodon", profile.MastodonAddress))
	return err
}
