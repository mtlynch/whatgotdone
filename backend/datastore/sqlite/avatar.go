package sqlite

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (d DB) GetAvatar(username types.Username) (io.Reader, error) {
	var avatar []byte
	if err := d.ctx.QueryRow(`
	SELECT
		avatar
	FROM
		user_profiles
	WHERE
		avatar IS NOT NULL AND
		username=?`, username).Scan(&avatar); err != nil {
		if err == sql.ErrNoRows {
			return nil, datastore.ErrAvatarNotFound{
				Username: username,
			}
		}
		return nil, err
	}

	if len(avatar) == 0 {
		return nil, errors.New("no avatar for user")
	}

	return bytes.NewBuffer(avatar), nil
}

func (d DB) InsertAvatar(username types.Username, avatarReader io.Reader, avatarWidth int) error {
	log.Printf("saving avatar to datastore for user %s", username)
	avatar, err := io.ReadAll(avatarReader)
	if err != nil {
		return err
	}
	if _, err := d.ctx.Exec(`
	UPDATE user_profiles
	SET
		avatar=?,
		avatar_width=?,
		avatar_last_modified=strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc')
	WHERE
		username=?`, avatar, avatarWidth, username); err != nil {
		return err
	}
	return nil
}

func (d DB) DeleteAvatar(username types.Username) error {
	log.Printf("deleting avatar from datastore for user %s", username)
	if _, err := d.ctx.Exec(`
	UPDATE user_profiles
	SET
		avatar=NULL,
		avatar_width=NULL,
		avatar_last_modified=NULL
	WHERE
		username=?`); err != nil {
		return err
	}
	return nil
}
