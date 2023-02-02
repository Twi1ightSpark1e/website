package errors

import (
	"log"
	"net/http"
	"strings"
)

func AssertPath(path string, w http.ResponseWriter, r *http.Request, errlog *log.Logger) bool {
	if path == r.URL.Path {
		return true
	}

	WriteNotFoundError(w, r, errlog)
	return false
}

func AssertPathBeginning(path string, w http.ResponseWriter, r *http.Request, errlog *log.Logger) bool {
	if strings.Index(r.URL.Path, path) == 0 {
		return true
	}

	WriteNotFoundError(w, r, errlog)
	return false
}
