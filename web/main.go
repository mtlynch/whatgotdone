package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("./web/frontend/dist"))
	r := mux.NewRouter()
	r.PathPrefix("/js").Handler(fs)
	r.PathPrefix("/css").Handler(fs)
	r.HandleFunc("/entries", handlers.EntriesHandler)
	r.HandleFunc("/api/submit", handlers.SubmitHandler)
	r.PathPrefix("/").HandlerFunc(handlers.IndexHandler)
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
