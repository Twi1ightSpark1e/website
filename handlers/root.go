package handlers

import (
	"net"
	"net/http"

	configuration "github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
)

type cardLink struct {
	Title string
	Address string
}

type card struct {
	Title string
	Content string
	Links []cardLink
}

type rootPage struct {
	Host string
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
	Cards []card
	Error string
}

type rootHandler struct {
	logger log.Channels
}
func RootHandler(logger log.Channels) http.Handler {
	template.AssertExists("index", logger)
	return &rootHandler{logger}
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	breadcrumb := prepareBreadcrum(r)

	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	tplData := rootPage {
		Host: r.Host,
		Cards: getCards(remoteAddr),
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
	}
	if len(tplData.Breadcrumb) > 0 {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	err := template.Get("index").Execute(w, tplData)
	if err != nil {
		h.logger.Err.Print(err)
	}
}

func getCards(client net.IP) []card {
	config := configuration.Get()

	cards := []card {}
	for _, cardDescr := range config.RootContent {
		if !configuration.IsAllowedByACL(client, cardDescr.View) {
			continue
		}

		links := []cardLink {}
		for _, linkDescr := range cardDescr.Links {
			links = append(links, cardLink {
				Title: linkDescr.Title,
				Address: linkDescr.Address,
			})
		}

		cards = append(cards, card {
			Title: cardDescr.Title,
			Content: cardDescr.Description,
			Links: links,
		})
	}

	return cards
}

