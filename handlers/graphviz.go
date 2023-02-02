package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Twi1ightSpark1e/website/config"
	handlerrors "github.com/Twi1ightSpark1e/website/handlers/errors"
	"github.com/Twi1ightSpark1e/website/handlers/util"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type graphvizPage struct {
	util.BreadcrumbData
	Image string
	Timestamp string
}

type graphData struct {
	image bytes.Buffer
	timestamp int64
}

type graphvizHandler struct {
	logger log.Channels
	path string
	endpoint config.GraphvizEndpointStruct
	graph graphData
}
func GraphvizHandler(logger log.Channels, path string, endpoint config.GraphvizEndpointStruct) http.Handler {
	template.AssertExists("graphviz", logger)
	return &graphvizHandler{logger, path, endpoint, graphData{}}
}

func (h *graphvizHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tplData := graphvizPage {
		BreadcrumbData: util.PrepareBreadcrumb(r),
	}

	switch r.Method {
	case http.MethodPut:
		h.handlePUT(w, r)
		return
	case http.MethodDelete:
		h.handleDELETE(w, r)
		return
	case http.MethodGet:
		if !h.handleGET(w, r, &tplData) {
			return
		}
		err := util.MinifyTemplate("graphviz", tplData, w)
		if err != nil {
			h.logger.Err.Print(err)
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		handlerrors.WriteError(w, r, errors.New("Invalid request method"), h.logger.Err)
		return
	}
}

func (h *graphvizHandler) handlePUT(w http.ResponseWriter, r *http.Request) {
	remoteAddr := util.GetRemoteAddr(r)
	h.logger.Info.Printf("Client %s sent PUT request on '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.Edit) {
		handlerrors.WriteNotFoundError(w, r, h.logger.Err)
		return
	}

	if !handlerrors.AssertPath(h.path, w, r, h.logger.Err) {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handlerrors.WriteError(w, r, err, h.logger.Err)
		return
	}

	g := graphviz.New()
	graph, err := graphviz.ParseBytes(body)
	if err != nil {
		handlerrors.WriteError(w, r, err, h.logger.Err)
		return
	}

	h.performDecoration(g, graph)

	// Render graph
	var buffer bytes.Buffer
	if err = g.Render(graph, graphviz.SVG, &buffer); err != nil {
		handlerrors.WriteError(w, r, err, h.logger.Err)
		return
	}

	h.graph.image = buffer
	h.graph.timestamp = time.Now().Unix()

	w.Write([]byte("ok"))
}

func (h *graphvizHandler) handleDELETE(w http.ResponseWriter, r *http.Request) {
	remoteAddr := util.GetRemoteAddr(r)
	h.logger.Info.Printf("Client %s sent DELETE request on '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.Edit) {
		handlerrors.WriteNotFoundError(w, r, h.logger.Err)
		return
	}

	if !handlerrors.AssertPath(h.path, w, r, h.logger.Err) {
		return
	}

	h.graph = graphData{}
	w.Write([]byte("ok"))
}

func (h *graphvizHandler) handleGET(w http.ResponseWriter, r *http.Request, tpl *graphvizPage) bool {
	remoteAddr := util.GetRemoteAddr(r)
	h.logger.Info.Printf("Client %s sent GET request on '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		handlerrors.WriteNotFoundError(w, r, h.logger.Err)
		return false
	}

	if !handlerrors.AssertPath(h.path, w, r, h.logger.Err) {
		return false
	}

	tpl.Image = base64.StdEncoding.EncodeToString(h.graph.image.Bytes())

	if h.graph.timestamp == 0 {
		tpl.Timestamp = "not performed yet"
	} else {
		tpl.Timestamp = time.Unix(h.graph.timestamp, 0).String()
	}

	return true
}

func (h *graphvizHandler) performDecoration(g *graphviz.Graphviz, graph *cgraph.Graph) {
	if h.endpoint.Decoration == config.DecorationTinc {
		layout := graphviz.CIRCO
		if graph.NumberNodes() >= 10 {
			layout = graphviz.DOT
		}
		g.SetLayout(layout)

		graph.SetBackgroundColor("transparent")

		for node := graph.FirstNode(); node != nil; node = graph.NextNode(node) {
			if node.Get("style") == "filled" {
				node.SetFillColor(node.Get("color"))
			} else {
				node.SetStyle(cgraph.FilledNodeStyle).SetFillColor("#ffffff")
			}
		}
	}

	// Decoration is `none`, so nothing to do here
}
