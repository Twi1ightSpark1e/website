package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type breadcrumbItem struct {
	Title string
	Address string
}

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
