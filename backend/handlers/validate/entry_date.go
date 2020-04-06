package validate

import (
	"time"

	"github.com/mtlynch/whatgotdone/backend/dates"
)

// EntryDate validates that the given date is valid for a journal entry. To be
// valid, the date must be:
//
//  * In YYYY-MM-DD format
//  * A Friday
//  * After 2019-01-01
//  * Be no later than the nearest following Friday
func EntryDate(date string) bool {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	const whatGotDoneEpochYear = 2019
	if t.Year() < whatGotDoneEpochYear {
		return false
	}
	if t.Weekday() != time.Friday {
		return false
	}
	if t.After(dates.ThisFriday()) {
		return false
	}
	return true
}
