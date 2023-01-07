//go:build dev || staging

package handlers

import (
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore/sqlite"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// addDevRoutes adds debug routes that we only use during development or e2e
// tests.
func (s *defaultServer) addDevRoutes() {
	s.router.HandleFunc("/api/testing/db/populate-dummy-data", s.wipeDB()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/testing/db/wipe", s.wipeDB()).Methods(http.MethodGet)
}

type dummyUserData struct {
	profile     types.UserProfile
	preferences types.Preferences
	drafts      []types.JournalEntry
	entries     []types.JournalEntry
	following   []types.Username
	reactions   map[types.EntryDate][]types.Reaction
}

func (s defaultServer) populateDummyData() http.HandlerFunc {
	dummyData := map[types.Username]dummyUserData{
		types.Username("staging_jimmy"): {
			drafts: []types.JournalEntry{
				{
					Date:     types.EntryDate("2019-06-28"),
					Markdown: types.EntryContent("Today was a productive day. I created a test data manager for What Got Done!"),
				},
			},
			entries: []types.JournalEntry{
				{
					Date:     types.EntryDate("2019-06-21"),
					Markdown: types.EntryContent("Watched *The Terminator* and wondered whether robots will ever really have cool accents."),
				},
				{
					Date:     types.EntryDate("2019-06-28"),
					Markdown: types.EntryContent("Today was a productive day. I created a test data manager for What Got Done!"),
				},
			},
		},
		types.Username("leader_lenny"): {
			entries: []types.JournalEntry{
				{
					Date:     types.EntryDate("2012-12-03"),
					Markdown: types.EntryContent("It's good to be the leader, as other users love to follow me!"),
				},
				{
					Date:     types.EntryDate("2012-11-26"),
					Markdown: types.EntryContent("Ate some Hot Pockets and played the lottery."),
				},
			},
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for username, ud := range dummyData {
			if err := s.datastore.SetPreferences(username, types.Preferences{
				EntryTemplate: types.EntryContent(ud.preferences.EntryTemplate),
			}); err != nil {
				panic(err)
			}
			if err := s.datastore.SetUserProfile(username, types.UserProfile{
				AboutMarkdown:   ud.profile.AboutMarkdown,
				EmailAddress:    ud.profile.EmailAddress,
				TwitterHandle:   ud.profile.TwitterHandle,
				MastodonAddress: ud.profile.MastodonAddress,
			}); err != nil {
				panic(err)
			}
			for _, d := range ud.drafts {
				d.Author = username
				if err := s.datastore.InsertDraft(username, d); err != nil {
					panic(err)
				}
			}
			for _, e := range ud.entries {
				if err := s.datastore.InsertEntry(username, e); err != nil {
					panic(err)
				}
			}
			for _, leader := range ud.following {
				err := s.datastore.InsertFollow(leader, username)
				if err != nil {
					panic(err)
				}
			}
			for date, reactions := range ud.reactions {
				for _, r := range reactions {
					if err := s.datastore.AddReaction(username, date, r); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

// wipeDB wipes the database back to a freshly initialized state.
func (s defaultServer) wipeDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlStore, ok := s.datastore.(*sqlite.DB)
		if !ok {
			log.Fatalf("store is not SQLite, can't wipe database")
		}
		sqlStore.Clear()
	}
}
