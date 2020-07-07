package main

import (
	"github.com/chefe/unews/helper"
	"github.com/chefe/unews/srf"
	"github.com/chefe/unews/twentymin"
	"log"
	"net/http"
)

func main() {
	helper.RegisterFeed("/20min", twentymin.GetFeed)
	helper.RegisterFeed("/srf", srf.GetFeed)

	log.Print("Server is running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
