package handlers

import (
	"log"
	"net"
	"net/http"

	configuration "github.com/Twi1ightSpark1e/website/config"
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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	breadcrumb := prepareBreadcrum(r)

	remoteAddr := getRemoteAddr(r)
	tplData := rootPage {
		Host: r.Host,
		Cards: getCards(remoteAddr),
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
	}
	if len(tplData.Breadcrumb) > 0 {
		w.WriteHeader(http.StatusNotFound)
		tplData.Error = "Content not found"
	}

	err := template.Get("index").Execute(w, tplData)
	if err != nil {
		log.Print(err)
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

