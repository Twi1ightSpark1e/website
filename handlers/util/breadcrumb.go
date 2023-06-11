package util

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"
)

type BreadcrumbData struct {
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
	ThemeSwitch Theme
}

type breadcrumbItem struct {
	Title string
	Address string
}

func PrepareBreadcrumb(req *http.Request) BreadcrumbData {
	result := []breadcrumbItem {
		{
			Title: req.Host,
			Address: "/",
		},
	}

	items := lo.Filter(strings.Split(req.URL.Path, "/"), func (item string, index int) bool {
		return len(item) != 0
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

	return BreadcrumbData{
		Breadcrumb: result[:len(result) - 1],
		LastBreadcrumb: result[len(result) - 1].Title,
		ThemeSwitch: GetTheme(req),
	}
}
