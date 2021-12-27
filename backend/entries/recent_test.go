package entries

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestRecentObservesStartAndLimitParameters(t *testing.T) {
	entries := []types.JournalEntry{
		{Author: "bob", Date: "2019-05-10", LastModified: mustParseTime("2019-05-25T06:00:00.000Z"), Markdown: "Read the news today... Oh boy!"},
		{Author: "bob", Date: "2019-05-03", LastModified: mustParseTime("2019-05-16T00:00:00.000Z"), Markdown: "Took a nap and dreamed about chocolate"},
		{Author: "bob", Date: "2019-04-26", LastModified: mustParseTime("2019-05-25T00:00:00.000Z"), Markdown: "Read a book about the history of cheese"},
		{Author: "bob", Date: "2019-04-19", LastModified: mustParseTime("2019-05-17T12:00:00.000Z"), Markdown: "Saw a movie about French vanilla"},
		{Author: "bob", Date: "2019-04-12", LastModified: mustParseTime("2019-05-23T00:00:00.000Z"), Markdown: "Ate some crackers in a bathtub"},
		{Author: "bob", Date: "2019-04-05", LastModified: mustParseTime("2019-05-24T00:00:00.000Z"), Markdown: "Rode the bus and saw a movie about ghosts"},
	}
	ms := mock.MockDatastore{
		JournalEntries: entries,
		Usernames: []types.Username{
			"bob",
		},
	}
	r := defaultReader{
		store: &ms,
	}

	var tests = []struct {
		explanation     string
		start           int
		limit           int
		entriesExpected []types.JournalEntry
	}{
		{
			"observes valid start and limit values",
			1,
			3,
			[]types.JournalEntry{
				{Author: "bob", Date: "2019-05-03", LastModified: mustParseTime("2019-05-16T00:00:00.000Z"), Markdown: "Took a nap and dreamed about chocolate"},
				{Author: "bob", Date: "2019-04-26", LastModified: mustParseTime("2019-05-25T00:00:00.000Z"), Markdown: "Read a book about the history of cheese"},
				{Author: "bob", Date: "2019-04-19", LastModified: mustParseTime("2019-05-17T12:00:00.000Z"), Markdown: "Saw a movie about French vanilla"},
			},
		},
		{
			"accepts large ranges",
			0,
			500,
			[]types.JournalEntry{
				{Author: "bob", Date: "2019-05-10", LastModified: mustParseTime("2019-05-25T06:00:00.000Z"), Markdown: "Read the news today... Oh boy!"},
				{Author: "bob", Date: "2019-05-03", LastModified: mustParseTime("2019-05-16T00:00:00.000Z"), Markdown: "Took a nap and dreamed about chocolate"},
				{Author: "bob", Date: "2019-04-26", LastModified: mustParseTime("2019-05-25T00:00:00.000Z"), Markdown: "Read a book about the history of cheese"},
				{Author: "bob", Date: "2019-04-19", LastModified: mustParseTime("2019-05-17T12:00:00.000Z"), Markdown: "Saw a movie about French vanilla"},
				{Author: "bob", Date: "2019-04-12", LastModified: mustParseTime("2019-05-23T00:00:00.000Z"), Markdown: "Ate some crackers in a bathtub"},
				{Author: "bob", Date: "2019-04-05", LastModified: mustParseTime("2019-05-24T00:00:00.000Z"), Markdown: "Rode the bus and saw a movie about ghosts"},
			},
		},
		{
			"returns empty for start beyond size of total response",
			500,
			5,
			[]types.JournalEntry{},
		},
	}
	for _, tt := range tests {
		actual, err := r.Recent(tt.start, tt.limit)
		if err != nil {
			t.Fatalf("Failed to retrieve recent entries: %v", err)
		}
		if !reflect.DeepEqual(actual, tt.entriesExpected) {
			t.Fatalf("Unexpected response: got %+v want %+v", actual, tt.entriesExpected)
		}
	}
}

func TestRecentFailsWhenDatastoreFailsToRetrieveEntries(t *testing.T) {
	ms := mock.MockDatastore{
		Usernames: []types.Username{
			"bob",
		},
		ReadEntriesErr: errors.New("dummy error for MockDatastore.GetEntries()"),
	}
	r := defaultReader{
		store: &ms,
	}

	_, err := r.Recent(0, 20)
	if err == nil {
		t.Fatalf("Expected call to Recent to fail")
	}
}

func mustParseTime(ts string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", ts)
	if err != nil {
		panic(err)
	}
	return t
}
