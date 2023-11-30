package sqlite

import (
	"bytes"
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

	return bytes.NewBuffer(avatar), nil
}
