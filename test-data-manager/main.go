// +build dev staging

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	log.Print("Starting test-data-manager script")

	source := flag.String("source", "dev-data.yaml", "Path to JSON or YAML file with data to load")
	keepAlive := flag.Bool("keepAlive", false, "Stay alive after completing initialization")
	flag.Parse()

	log.Printf("source=%s", *source)

	data := loadFromFile(*source)

	mgr := NewManager(data)

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
