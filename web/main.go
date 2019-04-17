package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mtlynch/whatgotdone/handlers"
)

func main() {
	s := handlers.New()
	http.Handle("/", s.Router())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
