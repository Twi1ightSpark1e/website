package handlers

import (
	"net"
	"net/http"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
)

type cardsPage struct {
	breadcrumb
	Cards []config.CardStruct
}

type cardsHandler struct {
	logger log.Channels
	path string
	endpoint config.CardsEndpointStruct
}
func CardsHandler(logger log.Channels, path string, endpoint config.CardsEndpointStruct) http.Handler {
	template.AssertExists("cards", logger)
	return &cardsHandler{logger, path, endpoint}
}

func (h *cardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	if !assertPath(h.path, w, r, h.logger.Err) {
		return
	}

	tplData := cardsPage {
		Cards: h.getCards(remoteAddr),
		breadcrumb: prepareBreadcrum(r),
	}

	err := minifyTemplate("cards", tplData, w)
	if err != nil {
		h.logger.Err.Print(err)
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
