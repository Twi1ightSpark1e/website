package handlers

import (
	"bytes"
	"encoding/base64"
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
	Error string
}

type graphData struct {
	image bytes.Buffer
	timestamp int64
}

type graphvizHandler struct {
	logger log.Channels
	endpoint config.GraphvizEndpointStruct
	graph graphData
}
func GraphvizHandler(logger log.Channels, endpoint config.GraphvizEndpointStruct) http.Handler {
	template.AssertExists("graphviz", logger)
	return &graphvizHandler{logger, endpoint, graphData{}}
}

func (h *graphvizHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	breadcrumb := prepareBreadcrum(r)
	tplData := graphvizPage {
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
	}

	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	if len(tplData.Breadcrumb) > 1 {
		w.WriteHeader(http.StatusNotFound)
		tplData.Error = "Content not found"
		goto compile
	}

	switch r.Method {
	case "PUT":
		if err := h.HandlePUT(w, r); err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("ok"))
		}
		return
	case "GET":
		tplData.Image = base64.StdEncoding.EncodeToString(h.graph.image.Bytes())
	default:
		w.WriteHeader(http.StatusForbidden)
		tplData.Error = "Invalid request method"
	}

compile:
	err := template.Get("graphviz").Execute(w, tplData)
	if err != nil {
		h.logger.Err.Print(err)
	}
}

func (h *graphvizHandler) HandlePUT(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	g := graphviz.New().SetLayout("circo")
	graph, err := graphviz.ParseBytes(body)
	if err != nil {
		return err
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
		return err
	}

	h.graph.image = buffer
	h.graph.timestamp = time.Now().Unix()

	return nil
}
