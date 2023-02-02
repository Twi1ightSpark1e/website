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
	logger log.Channels
	cache *template.HTML
}
func CreateHandler(
	root http.FileSystem,
	path string,
	endpoint config.MarkdownEndpointStruct,
	logger log.Channels,
) http.Handler {
	tpl.AssertExists("markdown", logger)

	return &handler{root, path, endpoint, logger, nil}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := util.GetRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		errors.WriteNotFoundError(w, r, h.logger.Err)
		return
	}

	if !errors.AssertPath(h.path, w, r, h.logger.Err) {
		return
	}

	if h.cache == nil {
		html, err := h.render()
		if err != nil {
			errors.WriteError(w, r, err, h.logger.Err)
			return
		}
		h.cache = &html
		h.logger.Info.Printf("Cache rebuild triggered by '%s' request", r.URL.Path)
	}

	tpl := page{
		BreadcrumbData: util.PrepareBreadcrumb(r),
		Content: *h.cache,
	}
	err := util.MinifyTemplate("markdown", tpl, w)
	if err != nil {
		h.logger.Err.Print(err)
	}
}
