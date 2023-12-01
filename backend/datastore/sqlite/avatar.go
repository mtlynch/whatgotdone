package sqlite

import (
	"bytes"
	"database/sql"
	"errors"
	"io"

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
