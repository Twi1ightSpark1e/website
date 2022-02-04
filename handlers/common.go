package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

type breadcrumbItem struct {
	Title string
	Address string
}

var m = minify.New()

func prepareBreadcrum(req *http.Request) []breadcrumbItem {
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

	return result
}

func getRemoteAddr(req *http.Request) net.IP {
	if val, ok := req.Header["X-Real-Ip"]; ok {
		return net.ParseIP(val[0])
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return net.ParseIP(ip)
}

func InitializeMinify() {
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
	})
}

func minifyTemplate(t *template.Template, data interface{}, out io.Writer) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	if err := m.Minify("text/html", out, &buf); err != nil {
		return err
	}

	return nil
}
