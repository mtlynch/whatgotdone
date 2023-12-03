package sqlite

import (
	"bytes"
	"database/sql"
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
		avatars
	WHERE
		username=?`, username).Scan(&avatar); err != nil {
		if err == sql.ErrNoRows {
			return nil, datastore.ErrAvatarNotFound{
				Username: username,
			}
		}
		return nil, err
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
	INSERT OR REPLACE INTO avatars
	(
		username,
		avatar,
		width,
		last_modified
	)
	VALUES(
		?,
		?,
		?,
		strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc')
	)`, username, avatar, avatarWidth); err != nil {
		return err
	}

	return nil
}

func (d DB) DeleteAvatar(username types.Username) error {
	log.Printf("deleting avatar from datastore for user %s", username)
	if _, err := d.ctx.Exec(`
	DELETE FROM avatars
	WHERE username=?`, username); err != nil {
		return err
	}
	return nil
}
