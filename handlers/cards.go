package handlers

import (
	"net"
	"net/http"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/handlers/errors"
	"github.com/Twi1ightSpark1e/website/handlers/util"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
)

type cardsPage struct {
	util.BreadcrumbData
	Cards []config.CardStruct
}

type cardsHandler struct {
	path string
	endpoint config.CardsEndpointStruct
}
func CardsHandler(path string, endpoint config.CardsEndpointStruct) http.Handler {
	template.AssertExists("cards")
	return &cardsHandler{path, endpoint}
}

func (h *cardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := util.GetRemoteAddr(r)

	if !errors.AssertPath(h.path, w, r) {
		return
	}

	tplData := cardsPage {
		Cards: h.getCards(remoteAddr),
		BreadcrumbData: util.PrepareBreadcrumb(r),
	}

	err := util.MinifyTemplate("cards", tplData, w)
	if err != nil {
		log.Stderr().Print(err)
	}
}

func (h *cardsHandler) getCards(client net.IP) []config.CardStruct {
	cards := []config.CardStruct {}

	for _, cardDescr := range h.endpoint.Content {
		if config.IsAllowedByACL(client, cardDescr.View) {
			cards = append(cards, cardDescr)
		}
	}

	return cards
}
