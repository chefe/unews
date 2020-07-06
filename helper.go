package main

import (
	"fmt"
	"github.com/gorilla/feeds"
	"log"
	"net/http"
)

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type feedGetter = func() (*feeds.Feed, error)
type feedConverter = func(*feeds.Feed) (string, error)

func printFeed(w http.ResponseWriter, get feedGetter, convert feedConverter) {
	feed, err := get()
	logErr(err)
	resp, err := convert(feed)
	logErr(err)
	fmt.Fprintf(w, resp)
}

func printAtomFeed(w http.ResponseWriter, getFeed feedGetter) {
	printFeed(w, getFeed, func(feed *feeds.Feed) (string, error) {
		return feed.ToAtom()
	})
}

func printRssFeed(w http.ResponseWriter, getFeed feedGetter) {
	printFeed(w, getFeed, func(feed *feeds.Feed) (string, error) {
		return feed.ToRss()
	})
}

func printJSONFeed(w http.ResponseWriter, getFeed feedGetter) {
	printFeed(w, getFeed, func(feed *feeds.Feed) (string, error) {
		return feed.ToJSON()
	})
}
