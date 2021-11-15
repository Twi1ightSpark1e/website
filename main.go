package main

import (
	"log"
	"net/http"

	"github.com/Twi1ightSpark1e/website/handlers"
	"github.com/Twi1ightSpark1e/website/template"
)

func main() {
	template.Initialize()

	http.HandleFunc("/files/", handlers.FileindexHandler)
	http.HandleFunc("/files-test/", handlers.FileindexHandler)
	http.HandleFunc("/orel/", handlers.FileindexHandler)
	http.HandleFunc("/packages/", handlers.FileindexHandler)
	http.Handle("/static/", http.FileServer(http.Dir("/")))
	http.HandleFunc("/", handlers.RootHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
