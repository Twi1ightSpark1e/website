package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type errorPage struct {
	breadcrumb
	Error string
}

func assertPath(path string, w http.ResponseWriter, r *http.Request, errlog *log.Logger) bool {
	if path == r.URL.Path {
		return true
	}

	writeNotFoundError(w, r, errlog)
	return false
}

func assertPathBeginning(path string, w http.ResponseWriter, r *http.Request, errlog *log.Logger) bool {
	if strings.Index(r.URL.Path, path) == 0 {
		return true
	}

	writeNotFoundError(w, r, errlog)
	return false
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
	tpl := errorPage{
		breadcrumb: prepareBreadcrum(r),
		Error: message.Error(),
	}

	err := minifyTemplate("error", tpl, w)
	if err != nil {
		errlog.Print(err)
	}
}
