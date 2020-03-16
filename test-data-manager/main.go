// +build dev staging

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	"github.com/mtlynch/whatgotdone/backend/datastore"
)

func main() {
	log.Print("Starting bulk firestore tweak script")

	source := flag.String("source", "dev-data.yaml", "Path to YAML file with data to load")
	keepAlive := flag.Bool("keepAlive", false, "Stay alive after completing initialization")
	flag.Parse()

	log.Printf("source=%s", *source)

	mgr := NewManager(loadYaml(*source))
	waitForDatastore(mgr.datastore)

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

	log.Print("Exiting bulk firestore tweak script")
}

func waitForDatastore(ds datastore.Datastore) {
	retries := 10
	for i := 0; i < retries; i++ {
		log.Printf("contacting datastore - attempt #%d", i)
		_, err := ds.GetUserProfile("dummy_user")
		if _, ok := err.(datastore.UserProfileNotFoundError); ok {
			log.Print("successfully contacted datastore")
			return
		}
		log.Printf("datastore not yet ready (%v), retrying...", err)
		time.Sleep(100 * time.Millisecond)
	}
	panic("Failed to connect to datastore")
}
