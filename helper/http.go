package helper

import (
	"fmt"
	"github.com/gorilla/feeds"
	"log"
	"net/http"
)

type FeedGetter = func() (*feeds.Feed, error)

func RegisterFeed(path string, getFeed FeedGetter) {
	http.HandleFunc(path+".xml", func(w http.ResponseWriter, r *http.Request) {
		feed, err := getFeed()
		logErr(err)
		resp, err := feed.ToAtom()
		logErr(err)
		fmt.Fprintf(w, resp)
	})

	http.HandleFunc(path+".rss", func(w http.ResponseWriter, r *http.Request) {
		feed, err := getFeed()
		logErr(err)
		resp, err := feed.ToRss()
		logErr(err)
		fmt.Fprintf(w, resp)
	})

	http.HandleFunc(path+".json", func(w http.ResponseWriter, r *http.Request) {
		feed, err := getFeed()
		logErr(err)
		resp, err := feed.ToJSON()
		logErr(err)
		fmt.Fprintf(w, resp)
	})
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
