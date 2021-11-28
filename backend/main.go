package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	gorilla "github.com/gorilla/handlers"

	"github.com/mtlynch/whatgotdone/backend/handlers"
)

func main() {
	log.Print("Starting whatgotdone server")

	s := handlers.New()
	http.Handle("/", gorilla.ProxyHeaders(gorilla.LoggingHandler(os.Stdout, s.Router())))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("Listening on %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
