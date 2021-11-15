package handlers

import (
	"log"
	"net"
	"net/http"

	"github.com/Twi1ightSpark1e/website/template"
)

type cardLink struct {
	Title string
	Address string
}

type card struct {
	Title string
	Content string
	Links []cardLink
}

type rootPage struct {
	Host string
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
	Cards []card
	Error string
}

var baseCards []card
var hlebCards []card

var hlebV4net = net.IPNet {
	IP: net.ParseIP("10.41.0.0"),
	Mask: net.CIDRMask(16, 32),
}
var hlebV6net = net.IPNet {
	IP: net.ParseIP("fd00:41eb::"),
	Mask: net.CIDRMask(32, 128),
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	breadcrumb := prepareBreadcrum(r)
	tplData := rootPage {
		Host: r.Host,
		Cards: getBaseCards(),
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
	}
	if len(tplData.Breadcrumb) > 0 {
		w.WriteHeader(http.StatusNotFound)
		tplData.Error = "Content not found"
	}

	var remoteAddr net.IP
	log.Print(r)
	if val, ok := r.Header["X-Real-Ip"]; ok {
		remoteAddr = net.ParseIP(val[0])
	} else {
		remoteAddr = net.ParseIP(r.RemoteAddr)
	}
	if hlebV4net.Contains(remoteAddr) || hlebV6net.Contains(remoteAddr) {
		// tplData.Cards = append(tplData.Cards, getHlebCards()...)
	}

	err := template.Get("index").Execute(w, tplData)
	if err != nil {
		log.Print(err)
	}
}

func getBaseCards() []card {
	if len(baseCards) == 0 {
		baseCards = []card {
			{
				Title: "File storage",
				Content: "Title that speaks for itself",
				Links: []cardLink {
					{
						Title: "Link",
						Address: "files/",
					},
				},
			},
			{
				Title: "Gentoo binary host",
				Content: "Binhost that you cat freely use anytime you want to",
				Links: []cardLink {
					{
						Title: "Link",
						Address: "packages/",
					},
					{
						Title: "How to set up",
						Address: "https://wiki.gentoo.org/wiki/Binary_package_guide/en#Using_binary_packages",
					},
				},
			},
			// {
			// 	Title: "Grafana",
			// 	Content: "Grafana-based monitoring",
			// 	Links: []cardLink {
			// 		{
			// 			Title: "Server status",
			// 			Address: "grafana/d/qnuLkwOMk/server-status?orgId=2",
			// 		},
			// 	},
			// },
		}
	}
	return baseCards
}

func getHlebCards() []card {
	if len(hlebCards) == 0 {
		hlebCards = []card {
			{
				Title: "HlebMesh graph",
				Content: "Mesh status taken from tinc and visualised with graphviz",
				Links: []cardLink {
					{
						Title: "Link",
						Address: "hlebmesh-status/",
					},
				},
			},
		}
	}
	return hlebCards
}
