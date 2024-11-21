package sqlite

import (
	"database/sql"
	"log"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// InsertFollow adds a following relationship to the datastore.
func (d DB) InsertFollow(leader, follower types.Username) error {
	log.Printf("saving follow to datastore: %s follows %s", follower, leader)
	_, err := d.ctx.Exec(`
		INSERT OR REPLACE INTO follows(
				follower,
				leader,
				created)
		values(:follower, :leader, strftime('%Y-%m-%d %H:%M:%SZ', 'now', 'utc'))`,
		sql.Named("follower", follower),
		sql.Named("leader", leader))
	return err
}

// DeleteFollow removes a following relationship from the datastore.
func (d DB) DeleteFollow(leader, follower types.Username) error {
	log.Printf("deleting follow from datastore: %s stopped following %s", follower, leader)
	_, err := d.ctx.Exec(`
		DELETE FROM
				follows
		WHERE
				follower=:follower AND
				leader=:leader
		`,
		sql.Named("follower", follower),
		sql.Named("leader", leader))
	return err
}

// Following returns all the users the specified user is following.
func (d DB) Following(follower types.Username) ([]types.Username, error) {
	rows, err := d.ctx.Query(`
		SELECT
				leader
		FROM
				follows
		WHERE
				follower=:follower`,
		sql.Named("follower", follower))
	if err != nil {
		return []types.Username{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close SQLite rows: %v", err)
		}
	}()

	leaders := []types.Username{}
	for rows.Next() {
		var leader string
		err := rows.Scan(&leader)
		if err != nil {
			return []types.Username{}, err
		}
		leaders = append(leaders, types.Username(leader))
	}
	return leaders, nil
}
