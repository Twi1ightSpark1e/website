package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type graphvizPage struct {
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
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
	breadcrumb := prepareBreadcrum(r)
	tplData := graphvizPage {
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
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
		err := template.Get("graphviz").Execute(w, tplData)
		if err != nil {
			h.logger.Err.Print(err)
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		writeError(w, r, errors.New("Invalid request method"), h.logger.Err)
		return
	}
}

func (h *graphvizHandler) handlePUT(w http.ResponseWriter, r *http.Request) {
	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s sent PUT request on '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.Edit) {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	if !assertPath(h.path, w, r, h.logger.Err) {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, r, err, h.logger.Err)
		return
	}

	g := graphviz.New().SetLayout("circo")
	graph, err := graphviz.ParseBytes(body)
	if err != nil {
		writeError(w, r, err, h.logger.Err)
		return
	}

	// Styling graph
	graph.SetBackgroundColor("transparent")

	for node := graph.FirstNode(); node != nil; node = graph.NextNode(node) {
		if node.Get("style") == "filled" {
			node.SetFillColor(node.Get("color"))
		} else {
			node.SetStyle(cgraph.FilledNodeStyle).SetFillColor("#ffffff")
		}
	}

	// Render graph
	var buffer bytes.Buffer
	if err = g.Render(graph, graphviz.SVG, &buffer); err != nil {
		writeError(w, r, err, h.logger.Err)
		return
	}

	h.graph.image = buffer
	h.graph.timestamp = time.Now().Unix()

	w.Write([]byte("ok"))
}

func (h *graphvizHandler) handleDELETE(w http.ResponseWriter, r *http.Request) {
	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s sent DELETE request on '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.Edit) {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	if !assertPath(h.path, w, r, h.logger.Err) {
		return
	}

	h.graph = graphData{}
	w.Write([]byte("ok"))
}

func (h *graphvizHandler) handleGET(w http.ResponseWriter, r *http.Request, tpl *graphvizPage) bool {
	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s sent GET request on '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		writeNotFoundError(w, r, h.logger.Err)
		return false
	}

	if !assertPath(h.path, w, r, h.logger.Err) {
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
