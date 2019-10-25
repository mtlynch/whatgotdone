package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	muxHandlers "github.com/gorilla/handlers"

	"github.com/mtlynch/whatgotdone/backend/handlers"
)

func main() {
	log.Print("Starting whatgotdone server")

	datastoreAddr := flag.String("datastore", "", "Address of datastore to use (e.g., localhost:6379)")
	flag.Parse()

	s := handlers.New(*datastoreAddr)
	http.Handle("/", muxHandlers.LoggingHandler(os.Stdout, s.Router()))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Print("Options:")
	log.Printf("  datastore address: [%v]", *datastoreAddr)
	log.Printf("          HTTP port: [%v]", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
