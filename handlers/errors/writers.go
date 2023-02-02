package errors

import (
	"log"
	"net/http"

	"github.com/Twi1ightSpark1e/website/handlers/util"
	"github.com/pkg/errors"
)

type page struct {
	util.BreadcrumbData
	Error string
}

func WriteUnauthorizedError(w http.ResponseWriter, r *http.Request, errlog *log.Logger) {
	w.WriteHeader(http.StatusUnauthorized)
	WriteError(w, r, errors.New("Unauthorized"), errlog)
}

func WriteNotFoundError(w http.ResponseWriter, r *http.Request, errlog *log.Logger) {
	w.WriteHeader(http.StatusNotFound)
	WriteError(w, r, errors.New("Content not found"), errlog)
}

func WriteError(w http.ResponseWriter, r *http.Request, message error, errlog *log.Logger) {
	tpl := page{
		BreadcrumbData: util.PrepareBreadcrumb(r),
		Error: message.Error(),
	}

	err := util.MinifyTemplate("error", tpl, w)
	if err != nil {
		errlog.Print(err)
	}
}
