package sqlite

import (
	"bytes"
	"errors"
	"io"

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
		username=?`, username).Scan(&avatar); err != nil {
		return nil, err
	}

	if len(avatar) == 0 {
		return nil, errors.New("no avatar for user")
	}

	return bytes.NewBuffer(avatar), nil
}
