package dates

import "time"

// Returns the date of the soonest Friday.
func ThisFriday() time.Time {
	t := time.Now()
	for t.Weekday() != time.Friday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}
