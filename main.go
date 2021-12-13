package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	configuration "github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/handlers"
	"github.com/Twi1ightSpark1e/website/template"
)

func main() {
	template.Initialize()
	configuration.Initialize("config.yaml")
	config := configuration.Get()

	for entry := range config.Handlers.FileIndex.Endpoints {
		baseDir := http.Dir(config.Handlers.FileIndex.BasePath)
		endpoint := config.Handlers.FileIndex.Endpoints[entry]
		handler := handlers.FileindexHandler(baseDir, endpoint)

		path := fmt.Sprintf("/%s/", entry)
		http.Handle(path, handler)

		visiblePath := strings.TrimRight(config.Handlers.FileIndex.BasePath, "/") + path
		log.Printf("New 'fileindex' handler for '%s'", visiblePath)
	}

	http.HandleFunc("/", handlers.RootHandler)

	addr := fmt.Sprintf(":%d", config.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

