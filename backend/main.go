package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	gorilla "github.com/mtlynch/gorilla-handlers"

	"github.com/mtlynch/whatgotdone/backend/datastore/sqlite"
	"github.com/mtlynch/whatgotdone/backend/handlers"
)

func main() {
	log.Print("Starting whatgotdone server")

	dbPath := flag.String("db", "data/store.db", "path to database")
	flag.Parse()
	ensureDirExists(filepath.Dir(*dbPath))
	datastore := sqlite.New(*dbPath)

	plausibleDomain := os.Getenv("PLAUSIBLE_DOMAIN")

	h := gorilla.LoggingHandler(os.Stdout, handlers.New(datastore, plausibleDomain).Router())
	if os.Getenv("BEHIND_PROXY") != "" {
		h = gorilla.ProxyIPHeadersHandler(h)
	}
	http.Handle("/", h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6001"
	}
	log.Printf("listening on %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func ensureDirExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
