package main

import (
	"fmt"
	"net/http"
	"strings"

	configuration "github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/handlers"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
)

func main() {
	logger := log.New("Main")

	template.Initialize()
	configuration.Initialize("config.yaml")
	config := configuration.Get()

	fileindexLogger := log.New("FileindexHandler")
	for entry := range config.Handlers.FileIndex.Endpoints {
		baseDir := http.Dir(config.Handlers.FileIndex.BasePath)
		endpoint := config.Handlers.FileIndex.Endpoints[entry]
		handler := handlers.FileindexHandler(baseDir, endpoint, fileindexLogger)

		path := fmt.Sprintf("/%s/", entry)
		http.Handle(path, handler)

		visiblePath := strings.TrimRight(config.Handlers.FileIndex.BasePath, "/") + path
		logger.Info.Printf("New 'fileindex' handler for '%s'", visiblePath)
	}

	http.Handle("/", handlers.RootHandler(log.New("RootHandler")))

	addr := fmt.Sprintf(":%d", config.Port)
	logger.Info.Printf("Listening TCP on '%s'", addr)
	logger.Err.Fatal(http.ListenAndServe(addr, nil))
}
