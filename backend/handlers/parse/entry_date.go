package parse

import (
	"errors"
	"fmt"
	"time"

	"github.com/mtlynch/whatgotdone/backend/dates"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// EntryDate parses a raw date string into a EntryDate object, ensuring that the
// given date is valid for a journal entry.
//
// To be valid, the date must be:
//
//   - In YYYY-MM-DD format
//   - A Friday
//   - After 2019-01-01
//   - Be no later than the nearest following Friday
func EntryDate(date string) (types.EntryDate, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}
	const whatGotDoneEpochYear = 2019
	if t.Year() < whatGotDoneEpochYear {
		return "", fmt.Errorf("invalid date, must be after %d", whatGotDoneEpochYear)
	}
	if t.Weekday() != time.Friday {
		return "", errors.New("invalid date, must be a Friday")
	}
	if t.After(dates.ThisFriday()) {
		return "", errors.New("invalid date, too far in the future")
	}
	return types.EntryDate(date), nil
}
