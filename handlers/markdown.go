package handlers

import (
	"html/template"
	"io"
	"net/http"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	tpl "github.com/Twi1ightSpark1e/website/template"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type markdownPage struct {
	breadcrumb
	Content template.HTML
}

type markdownHandler struct {
	root http.FileSystem
	path string
	endpoint config.MarkdownEndpointStruct
	logger log.Channels
	cache *template.HTML
}
func MarkdownHandler(
	root http.FileSystem,
	path string,
	endpoint config.MarkdownEndpointStruct,
	logger log.Channels,
) http.Handler {
	tpl.AssertExists("markdown", logger)

	return &markdownHandler{root, path, endpoint, logger, nil}
}

func (h *markdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	if !assertPath(h.path, w, r, h.logger.Err) {
		return
	}

	if h.cache == nil {
		html, err := h.renderMarkdown()
		if err != nil {
			writeError(w, r, err, h.logger.Err)
			return
		}
		h.cache = &html
		h.logger.Info.Printf("Cache rebuild triggered by '%s' request", r.URL.Path)
	}

	tpl := markdownPage{
		breadcrumb: prepareBreadcrum(r),
		Content: *h.cache,
	}
	err := minifyTemplate("markdown", tpl, w)
	if err != nil {
		h.logger.Err.Print(err)
	}
}

func (h *markdownHandler) renderMarkdown() (template.HTML, error) {
	file, err := h.root.Open(h.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	md, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return renderMarkdown(md), nil
}

func renderMarkdown(content []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.Attributes
	parser := parser.NewWithExtensions(extensions)
	html := template.HTML(markdown.ToHTML(content, parser, nil))

	return html
}
