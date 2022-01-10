package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Twi1ightSpark1e/website/template"
)

type errorPage struct {
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
	Error string
}

func writeUnauthorizedError(w http.ResponseWriter, r *http.Request, errlog *log.Logger) {
	w.WriteHeader(http.StatusUnauthorized)
	writeError(w, r, errors.New("Unauthorized"), errlog)
}

func writeNotFoundError(w http.ResponseWriter, r *http.Request, errlog *log.Logger) {
	w.WriteHeader(http.StatusNotFound)
	writeError(w, r, errors.New("Content not found"), errlog)
}

func writeError(w http.ResponseWriter, r *http.Request, message error, errlog *log.Logger) {
	breadcrumb := prepareBreadcrum(r)
	tpl := errorPage{
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
		Error: message.Error(),
	}

	err := template.Get("error").Execute(w, tpl)
	if err != nil {
		errlog.Print(err)
	}
}
