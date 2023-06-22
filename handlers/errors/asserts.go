package errors

import (
	"net/http"
	"strings"
)

func AssertPath(path string, w http.ResponseWriter, r *http.Request) bool {
	if path == r.URL.Path {
		return true
	}

	WriteNotFoundError(w, r)
	return false
}

func AssertPathBeginning(path string, w http.ResponseWriter, r *http.Request) bool {
	if strings.Index(r.URL.Path, path) == 0 {
		return true
	}

	WriteNotFoundError(w, r)
	return false
}
