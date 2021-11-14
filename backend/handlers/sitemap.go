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
		w.Write(sm.XMLContent())
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
	users, err := ds.Users()
	if err != nil {
		return
	}
	for _, u := range users {
		sm.Add(stm.URL{{"loc", fmt.Sprintf("/%s", u)}})
		entries, err := ds.ReadEntries(datastore.EntryFilter{
			ByUsers: []types.Username{u},
		})
		if err != nil {
			log.Printf("error getting entries for %s: %v", u, err)
			continue
		}
		for _, e := range entries {
			sm.Add(stm.URL{{"loc", fmt.Sprintf("/%s/%s", u, e.Date)}})
		}
	}
}
