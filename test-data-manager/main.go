//go:build dev || staging

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/mtlynch/whatgotdone/backend/datastore/sqlite"
)

func main() {
	log.Print("Starting test-data-manager script")

	source := flag.String("source", "dev-data.yaml", "Path to JSON or YAML file with data to load")
	dbPath := flag.String("db", "data/store.db", "path to database")
	keepAlive := flag.Bool("keepAlive", false, "Stay alive after completing initialization")
	flag.Parse()

	log.Printf("source=%s", *source)

	mgr := NewManager(loadFromFile(*source), sqlite.New(*dbPath))

	if *keepAlive {
		s := NewServer(mgr)
		http.Handle("/", handlers.LoggingHandler(os.Stdout, s.router))
		s.router.HandleFunc("/reset", s.resetPost()).Methods(http.MethodPost)

		port := os.Getenv("PORT")
		if port == "" {
			port = "5200"
		}
		log.Printf("Listening on %s", port)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	} else {
		mgr.Reset()
	}

	log.Print("Exiting test-data-manager script")
}
