package entries

import (
	"reflect"
	"testing"

	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestRecentSortsByDateThenByModifedTimeInDescendingOrder(t *testing.T) {
	entries := []types.JournalEntry{
		{Date: "2019-05-24", LastModified: "2019-05-24T00:00:00.000Z", Markdown: "Rode the bus and saw a movie about ghosts"},
		{Date: "2019-05-24", LastModified: "2019-05-23T00:00:00.000Z", Markdown: "Ate some crackers in a bathtub"},
		{Date: "2019-05-17", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
		{Date: "2019-05-24", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
		{Date: "2019-05-24", LastModified: "2019-05-25T22:00:00.000Z", Markdown: "Read a pamphlet from The Cat Society"},
		{Date: "2019-05-24", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
		{Date: "2019-05-17", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
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

	actual, err := r.Recent(0, 20)
	if err != nil {
		t.Fatalf("Failed to retrieve recent entries: %v", err)
	}

	// For simplicity of the test, all users have username "bob," but in practice,
	// these updates would come from different users.
	expected := []RecentEntry{
		{Author: "bob", Date: "2019-05-24", LastModified: "2019-05-25T22:00:00.000Z", Markdown: "Read a pamphlet from The Cat Society"},
		{Author: "bob", Date: "2019-05-24", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
		{Author: "bob", Date: "2019-05-24", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
		{Author: "bob", Date: "2019-05-24", LastModified: "2019-05-24T00:00:00.000Z", Markdown: "Rode the bus and saw a movie about ghosts"},
		{Author: "bob", Date: "2019-05-24", LastModified: "2019-05-23T00:00:00.000Z", Markdown: "Ate some crackers in a bathtub"},
		{Author: "bob", Date: "2019-05-17", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
		{Author: "bob", Date: "2019-05-17", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected response: got %v want %v", actual, expected)
	}
}

func TestRecentAlwaysPlacesNewDatesAheadOfOldDates(t *testing.T) {
	entries := []types.JournalEntry{
		{Date: "2019-05-17", LastModified: "2019-09-28T12:00:00.000Z", Markdown: "Made a hat out of donuts from the cloud in the sky"},
		{Date: "2019-09-20", LastModified: "2019-09-25T00:00:00.000Z", Markdown: "High-fived a platypus when the apple hits the pie."},
		{Date: "2019-09-06", LastModified: "2019-09-22T00:00:00.000Z", Markdown: "Ate an apple in a single bite of chocolate"},
		{Date: "2019-09-20", LastModified: "2019-09-20T00:00:00.000Z", Markdown: "Attended an Indie Hackers meetup"},
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

	actual, err := r.Recent(0, 20)
	if err != nil {
		t.Fatalf("Failed to retrieve recent entries: %v", err)
	}

	// For simplicity of the test, all users have username "bob," but in practice,
	// these updates would come from different users.
	expected := []RecentEntry{
		{Author: "bob", Date: "2019-09-20", LastModified: "2019-09-25T00:00:00.000Z", Markdown: "High-fived a platypus when the apple hits the pie."},
		{Author: "bob", Date: "2019-09-20", LastModified: "2019-09-20T00:00:00.000Z", Markdown: "Attended an Indie Hackers meetup"},
		{Author: "bob", Date: "2019-09-06", LastModified: "2019-09-22T00:00:00.000Z", Markdown: "Ate an apple in a single bite of chocolate"},
		{Author: "bob", Date: "2019-05-17", LastModified: "2019-09-28T12:00:00.000Z", Markdown: "Made a hat out of donuts from the cloud in the sky"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Unexpected response: got %v want %v", actual, expected)
	}
}

func TestRecentObservesStartAndLimitParameters(t *testing.T) {
	entries := []types.JournalEntry{
		{Date: "2019-05-10", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
		{Date: "2019-05-03", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
		{Date: "2019-04-26", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
		{Date: "2019-04-19", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
		{Date: "2019-04-12", LastModified: "2019-05-23T00:00:00.000Z", Markdown: "Ate some crackers in a bathtub"},
		{Date: "2019-04-05", LastModified: "2019-05-24T00:00:00.000Z", Markdown: "Rode the bus and saw a movie about ghosts"},
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
		entriesExpected []RecentEntry
	}{
		{
			"observes valid start and limit values",
			1,
			3,
			[]RecentEntry{
				{Author: "bob", Date: "2019-05-03", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
				{Author: "bob", Date: "2019-04-26", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
				{Author: "bob", Date: "2019-04-19", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
			},
		},
		{
			"accepts large ranges",
			0,
			500,
			[]RecentEntry{
				{Author: "bob", Date: "2019-05-10", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
				{Author: "bob", Date: "2019-05-03", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
				{Author: "bob", Date: "2019-04-26", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
				{Author: "bob", Date: "2019-04-19", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
				{Author: "bob", Date: "2019-04-12", LastModified: "2019-05-23T00:00:00.000Z", Markdown: "Ate some crackers in a bathtub"},
				{Author: "bob", Date: "2019-04-05", LastModified: "2019-05-24T00:00:00.000Z", Markdown: "Rode the bus and saw a movie about ghosts"},
			},
		},
		{
			"returns empty for start beyond size of total response",
			500,
			5,
			[]RecentEntry{},
		},
	}
	for _, tt := range tests {
		actual, err := r.Recent(tt.start, tt.limit)
		if err != nil {
			t.Fatalf("Failed to retrieve recent entries: %v", err)
		}
		if !reflect.DeepEqual(actual, tt.entriesExpected) {
			t.Fatalf("Unexpected response: got %v want %v", actual, tt.entriesExpected)
		}
	}
}
