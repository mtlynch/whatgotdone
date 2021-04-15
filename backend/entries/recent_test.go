package entries

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mtlynch/whatgotdone/backend/types"
)

type mockStore struct {
	journalEntries []types.JournalEntry
	users          []string
}

func (ms mockStore) Users() ([]string, error) {
	return ms.users, nil
}

func (ms mockStore) GetEntries(username string) ([]types.JournalEntry, error) {
	return ms.journalEntries, nil
}

func (ms mockStore) Following(follower string) ([]string, error) {
	return []string{}, errors.New("not implemented")
}

func (ms mockStore) Close() error {
	return nil
}

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
	ms := mockStore{
		journalEntries: entries,
		users: []string{
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

func TestRecentEntriesHandlerAlwaysPlacesNewDatesAheadOfOldDates(t *testing.T) {
	entries := []types.JournalEntry{
		{Date: "2019-05-17", LastModified: "2019-09-28T12:00:00.000Z", Markdown: "Made a hat out of donuts from the cloud in the sky"},
		{Date: "2019-09-20", LastModified: "2019-09-25T00:00:00.000Z", Markdown: "High-fived a platypus when the apple hits the pie."},
		{Date: "2019-09-06", LastModified: "2019-09-22T00:00:00.000Z", Markdown: "Ate an apple in a single bite of chocolate"},
		{Date: "2019-09-20", LastModified: "2019-09-20T00:00:00.000Z", Markdown: "Attended an Indie Hackers meetup"},
	}
	ms := mockStore{
		journalEntries: entries,
		users: []string{
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
