package markdown

import (
	"html/template"
	"net/http"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/handlers/errors"
	"github.com/Twi1ightSpark1e/website/handlers/util"
	"github.com/Twi1ightSpark1e/website/log"
	tpl "github.com/Twi1ightSpark1e/website/template"
)

type page struct {
	util.BreadcrumbData
	Content template.HTML
}

type handler struct {
	root http.FileSystem
	path string
	endpoint config.MarkdownEndpointStruct
	cache *template.HTML
}
func CreateHandler(
	root http.FileSystem,
	path string,
	endpoint config.MarkdownEndpointStruct,
) http.Handler {
	tpl.AssertExists("markdown")

	return &handler{root, path, endpoint, nil}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := util.GetRemoteAddr(r)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		errors.WriteNotFoundError(w, r)
		return
	}

	if !errors.AssertPath(h.path, w, r) {
		return
	}

	if h.cache == nil {
		html, err := h.render()
		if err != nil {
			errors.WriteError(w, r, err)
			return
		}
		h.cache = &html
	}

	tpl := page{
		BreadcrumbData: util.PrepareBreadcrumb(r),
		Content: *h.cache,
	}
	err := util.MinifyTemplate("markdown", tpl, w)
	if err != nil {
		log.Stderr().Print(err)
	}
}
