package main

import (
	"net/http"

	"github.com/Twi1ightSpark1e/website/handlers"
	"github.com/Twi1ightSpark1e/website/template"

	"unit.nginx.org/go"
)


func main() {
	template.Initialize()

	http.HandleFunc("/files/", handlers.FileindexHandler)
	http.HandleFunc("/files-test/", handlers.FileindexHandler)
	http.HandleFunc("/packages/", handlers.FileindexHandler)
	http.HandleFunc("/", handlers.RootHandler)

	unit.ListenAndServe(":8081", nil)
}

