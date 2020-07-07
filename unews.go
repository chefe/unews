package main

import (
	"github.com/chefe/unews/srf"
	"github.com/chefe/unews/twentymin"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/20min.xml", func(w http.ResponseWriter, r *http.Request) {
		printAtomFeed(w, twentymin.GetFeed)
	})

	http.HandleFunc("/20min.rss", func(w http.ResponseWriter, r *http.Request) {
		printRssFeed(w, twentymin.GetFeed)
	})

	http.HandleFunc("/20min.json", func(w http.ResponseWriter, r *http.Request) {
		printJSONFeed(w, twentymin.GetFeed)
	})

	http.HandleFunc("/srf.xml", func(w http.ResponseWriter, r *http.Request) {
		printAtomFeed(w, srf.GetFeed)
	})

	http.HandleFunc("/srf.rss", func(w http.ResponseWriter, r *http.Request) {
		printRssFeed(w, srf.GetFeed)
	})

	http.HandleFunc("/srf.json", func(w http.ResponseWriter, r *http.Request) {
		printJSONFeed(w, srf.GetFeed)
	})

	log.Print("Server is running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
