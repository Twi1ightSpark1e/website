package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	getopt "github.com/pborman/getopt/v2"

	configuration "github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/handlers"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
)

var (
	configPath = "config.yaml"
	showHelp = false
)

func main() {
	logger := log.New("Main")

	initialize()
	if showHelp {
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	}

	template.Initialize()
	configuration.Initialize(configPath)
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

func initialize() {
	getopt.FlagLong(&showHelp, "help", 'h', "Show usage information and exit.")
	getopt.FlagLong(&configPath, "config", 'c', "Path to configuration file.")

	getopt.Parse()
}
