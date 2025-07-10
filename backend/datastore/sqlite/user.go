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
		aboutMarkdown     string
		email             string
		twitter           string
		mastodon          string
		forwardingAddress sql.NullString
	)
	err := d.ctx.QueryRow(`
		SELECT
				about_markdown,
				email,
				twitter,
				mastodon,
				forwarding_address
		FROM
				user_profiles
		WHERE
				username = :username`,
		sql.Named("username", username)).Scan(&aboutMarkdown, &email, &twitter, &mastodon, &forwardingAddress)

	if err == sql.ErrNoRows {
		return types.UserProfile{}, datastore.UserProfileNotFoundError{
			Username: username,
		}
	} else if err != nil {
		return types.UserProfile{}, err
	}

	var forwardingAddr types.ForwardingAddress
	if forwardingAddress.Valid {
		forwardingAddr = types.ForwardingAddress(forwardingAddress.String)
	}

	return types.UserProfile{
		AboutMarkdown:     types.UserBio(aboutMarkdown),
		EmailAddress:      types.EmailAddress(email),
		TwitterHandle:     types.TwitterHandle(twitter),
		MastodonAddress:   types.MastodonAddress(mastodon),
		ForwardingAddress: forwardingAddr,
	}, nil
}

// SetUserProfile updates the given user's profile.
func (d DB) SetUserProfile(username types.Username, profile types.UserProfile) error {
	log.Printf("saving user profile to datastore: %s -> %+v", username, profile)

	var forwardingAddress sql.NullString
	if profile.ForwardingAddress != "" {
		forwardingAddress = sql.NullString{String: string(profile.ForwardingAddress), Valid: true}
	}

	_, err := d.ctx.Exec(`
		INSERT OR REPLACE INTO user_profiles(
				username,
				about_markdown,
				email,
				twitter,
				mastodon,
				forwarding_address)
		values(:username, :about_markdown, :email, :twitter, :mastodon, :forwarding_address)`,
		sql.Named("username", username),
		sql.Named("about_markdown", profile.AboutMarkdown),
		sql.Named("email", profile.EmailAddress),
		sql.Named("twitter", profile.TwitterHandle),
		sql.Named("mastodon", profile.MastodonAddress),
		sql.Named("forwarding_address", forwardingAddress))
	return err
}

// GetForwardingAddress retrieves the user's forwarding address.
func (d DB) GetForwardingAddress(username types.Username) (types.ForwardingAddress, error) {
	var forwardingAddress sql.NullString
	err := d.ctx.QueryRow(`
		SELECT forwarding_address
		FROM user_profiles
		WHERE username = :username`,
		sql.Named("username", username)).Scan(&forwardingAddress)

	if err == sql.ErrNoRows {
		return types.ForwardingAddress(""), datastore.ForwardingAddressNotFoundError{
			Username: username,
		}
	} else if err != nil {
		return types.ForwardingAddress(""), err
	}

	if !forwardingAddress.Valid || forwardingAddress.String == "" {
		return types.ForwardingAddress(""), datastore.ForwardingAddressNotFoundError{
			Username: username,
		}
	}

	return types.ForwardingAddress(forwardingAddress.String), nil
}

// SetForwardingAddress saves the user's forwarding address.
func (d DB) SetForwardingAddress(username types.Username, address types.ForwardingAddress) error {
	log.Printf("saving forwarding address to datastore: %s -> %s", username, address)

	// First ensure the user has a profile row
	_, err := d.ctx.Exec(`
		INSERT OR IGNORE INTO user_profiles(username, about_markdown, email, twitter, mastodon)
		VALUES(:username, '', '', '', '')`,
		sql.Named("username", username))
	if err != nil {
		return err
	}

	// Then update the forwarding address
	_, err = d.ctx.Exec(`
		UPDATE user_profiles
		SET forwarding_address = :forwarding_address
		WHERE username = :username`,
		sql.Named("username", username),
		sql.Named("forwarding_address", string(address)))
	return err
}

// DeleteForwardingAddress removes the user's forwarding address.
func (d DB) DeleteForwardingAddress(username types.Username) error {
	log.Printf("deleting forwarding address from datastore: %s", username)
	_, err := d.ctx.Exec(`
		UPDATE user_profiles
		SET forwarding_address = NULL
		WHERE username = :username`,
		sql.Named("username", username))
	return err
}
