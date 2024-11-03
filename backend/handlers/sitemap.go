package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s defaultServer) sitemapGet() http.HandlerFunc {
	var sm *stm.Sitemap = nil
	return func(w http.ResponseWriter, r *http.Request) {
		// Lazy-load the sitemap on first request.
		if sm == nil {
			sm = buildSitemap(s.datastore)
		}
		if _, err := w.Write(sm.XMLContent()); err != nil {
			log.Fatalf("failed to write sitemap: %v", err)
		}
	}
}

func buildSitemap(ds datastore.Datastore) *stm.Sitemap {
	sm := stm.NewSitemap(1)
	sm.SetDefaultHost("https://whatgotdone.com")

	sm.Create()
	sm.Add(stm.URL{{"loc", "/"}, {"changefreq", "daily"}})
	sm.Add(stm.URL{{"loc", "/about"}, {"changefreq", "daily"}})
	sm.Add(stm.URL{{"loc", "/recent"}, {"changefreq", "daily"}})
	sm.Add(stm.URL{{"loc", "/privacy-policy"}, {"changefreq", "daily"}})
	addUsersAndEntries(sm, ds)

	return sm
}

func addUsersAndEntries(sm *stm.Sitemap, ds datastore.Datastore) {
	entries, err := ds.ReadEntries(datastore.EntryFilter{})
	if err != nil {
		return
	}
	users := map[types.Username]bool{}
	// Add URLs for journal entries.
	for _, e := range entries {
		sm.Add(stm.URL{{"loc", fmt.Sprintf("/%s/%s", e.Author, e.Date)}})
		users[e.Author] = true
	}
	// Add URLs for user profiles.
	for u := range users {
		sm.Add(stm.URL{{"loc", fmt.Sprintf("/%s", u)}})
	}

}
