package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/Twi1ightSpark1e/website/template"
)

type breadcrumb struct {
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
}

type breadcrumbItem struct {
	Title string
	Address string
}

func prepareBreadcrum(req *http.Request) breadcrumb {
	result := []breadcrumbItem {
		{
			Title: req.Host,
			Address: "/",
		},
	}

	items := filterStr(strings.Split(req.URL.Path, "/"), func (item *string) bool {
		return len(*item) != 0
	})
	for idx, item := range items {
		if len(item) == 0 {
			continue
		}

		address := fmt.Sprintf("/%s/", strings.Join(items[:idx + 1], "/"))
		result = append(result, breadcrumbItem {
			Title: item,
			Address: address,
		})
	}

	return breadcrumb{
		Breadcrumb: result[:len(result) - 1],
		LastBreadcrumb: result[len(result) - 1].Title,
	}
}

func getRemoteAddr(req *http.Request) net.IP {
	if val, ok := req.Header["X-Real-Ip"]; ok {
		return net.ParseIP(val[0])
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return net.ParseIP(ip)
}

var m = minify.New()

func InitializeMinify() {
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
	})
}

func minifyTemplate(name string, data interface{}, out io.Writer) error {
	var buf bytes.Buffer
	if err := template.Execute(name, data, &buf); err != nil {
		return err
	}

	if err := m.Minify("text/html", out, &buf); err != nil {
		return err
	}

	return nil
}
