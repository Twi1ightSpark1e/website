package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

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

	configuration.Initialize(configPath)
	config := configuration.Get()
	template.Initialize()
	handlers.InitializeMinify()

	baseDir := http.Dir(config.Paths.Base)

	fileindexLogger := log.New("FileindexHandler")
	for entry, endpoint := range config.Handlers.FileIndex.Endpoints {
		path := handlerPath(entry)
		handler := handlers.FileindexHandler(baseDir, path, endpoint, fileindexLogger)
		http.Handle(path, handler)

		logger.Info.Printf("New 'fileindex' handler for '%s'", path)
	}

	graphvizLogger := log.New("GraphvizLogger")
	for entry, endpoint := range config.Handlers.Graphviz.Endpoints {
		path := handlerPath(entry)
		http.Handle(path, handlers.GraphvizHandler(graphvizLogger, path, endpoint))

		logger.Info.Printf("New 'graphviz' handler for '%s'", path)
	}

	webhookLogger := log.New("WebhookLogger")
	for entry, endpoint := range config.Handlers.Webhook.Endpoints {
		path := handlerPath(entry)
		http.Handle(path, handlers.WebhookHandler(webhookLogger, path, endpoint))

		logger.Info.Printf("New 'webhook' handler for '%s'", path)
	}

	cardsLogger := log.New("CardsLogger")
	for entry, endpoint := range config.Handlers.Cards.Endpoints {
		path := handlerPath(entry)
		http.Handle(path, handlers.CardsHandler(cardsLogger, path, endpoint))

		logger.Info.Printf("New 'cards' handler for '%s'", path)
	}

	markdownLogger := log.New("MarkdownLogger")
	for entry, endpoint := range config.Handlers.Markdown.Endpoints {
		path := handlerPath(entry)
		path = path[:len(path)-1]
		handler := handlers.MarkdownHandler(baseDir, path, endpoint, markdownLogger)
		http.Handle(path, handler)

		logger.Info.Printf("New 'markdown' handler for '%s'", path)
	}

	var wg sync.WaitGroup
	for _, addr := range config.Listen {
		wg.Add(1)

		logger.Info.Printf("Listening TCP on '%s'", addr)
		go func(addr string) {
			defer wg.Done()
			logger.Err.Fatal(http.ListenAndServe(addr, nil))
		}(addr)
	}
	wg.Wait()
}

func handlerPath(name string) string {
	if name != "index" {
		return fmt.Sprintf("/%s/", name)
	}
	return "/"
}

func initialize() {
	getopt.FlagLong(&showHelp, "help", 'h', "Show usage information and exit.")
	getopt.FlagLong(&configPath, "config", 'c', "Path to configuration file.")

	getopt.Parse()
}
