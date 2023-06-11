package fileindex

import (
	"fmt"
	"net/http"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/handlers/errors"
	"github.com/Twi1ightSpark1e/website/handlers/markdown"
	"github.com/Twi1ightSpark1e/website/handlers/util"
	"github.com/Twi1ightSpark1e/website/log"
	tpl "github.com/Twi1ightSpark1e/website/template"
)

type uploader func(w http.ResponseWriter, r *http.Request, params searchParams) (bool, error)

type handler struct {
	root http.FileSystem
	path string
	endpoint config.FileindexHandlerEndpointStruct
	logger log.Channels
	uploaders map[string]uploader
}

func CreateHandler(
	root http.FileSystem,
	path string,
	endpoint config.FileindexHandlerEndpointStruct,
	logger log.Channels,
) http.Handler {
	tpl.AssertExists("fileindex", logger)

	h := &handler{root, path, endpoint, logger, map[string]uploader{}}
	h.uploaders = map[string]uploader {
		"tar": h.uploadTar,
		"gz": h.uploadGz,
		"zst": h.uploadZst,
	}

	return h
}

type preservedParam struct {
	Key string
	Value string
}

type page struct {
	util.BreadcrumbData
	markdown.InlineMarkdown
	searchParams

	PreservedParams []preservedParam
	URL string
	AllowUpload bool
	List []fileEntry
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := util.GetRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	allowUpload := config.IsAllowedByACL(remoteAddr, h.endpoint.Upload)
	allowPost := r.Method == http.MethodPost && allowUpload
	allowView := r.Method != http.MethodPost && config.IsAllowedByACL(remoteAddr, h.endpoint.View)
	if !allowPost && !allowView {
		errors.WriteNotFoundError(w, r, h.logger.Err)
		return
	}

	if !errors.AssertPathBeginning(h.path, w, r, h.logger.Err) {
		return
	}

	pageData := page {
		BreadcrumbData: util.PrepareBreadcrumb(r),
		AllowUpload: allowUpload,
		URL: r.URL.Path,
		PreservedParams: h.preserveGetParams(r),
		searchParams: searchParams{
			FindQuery: r.URL.Query().Get("query"),
			FindMatchCase: r.URL.Query().Get("matchcase") == "on",
			FindRegex: r.URL.Query().Get("regex") == "on",
		},
	}
	hasQuery := len(pageData.searchParams.FindQuery) > 0
	pageData.AllowUpload = pageData.AllowUpload && !hasQuery

	if !config.Authenticate(r, h.endpoint.Auth) {
		authHeader := fmt.Sprintf(`Basic realm="Authentication required to use %s"`, pageData.LastBreadcrumb)
		w.Header().Set("WWW-Authenticate", authHeader)
		errors.WriteUnauthorizedError(w, r, h.logger.Err)
		return
	}

	if recv, err := h.recvFile(w, r); recv {
		return
	} else if err != nil {
		errors.WriteError(w, r, err, h.logger.Err)
		return
	}

	if sent, err := h.sendFile(w, r, pageData.searchParams); sent {
		return
	} else if err != nil {
		errors.WriteNotFoundError(w, r, h.logger.Err)
		return
	}

	if list, err := h.prepareFileList(r.URL.Path, remoteAddr, pageData.searchParams); err != nil {
		errors.WriteError(w, r, err, h.logger.Err)
		return
	} else {
		pageData.List = list

		show, name := h.showMarkdown(list)
		ptype := config.PreviewNone
		if show && !hasQuery {
			ptype = h.endpoint.Preview
			path := fmt.Sprintf("%s/%s", r.URL.Path, name)
			file, _ := h.root.Open(path)
			pageData.InlineMarkdown = markdown.PrepareInline(ptype, file)
			file.Close()
		}
	}

	err := util.MinifyTemplate("fileindex", pageData, w)
	if err != nil {
		h.logger.Err.Print(err)
	}
}

func (h *handler) preserveGetParams(r *http.Request) []preservedParam {
	result := make([]preservedParam, 0)
	for key, value := range r.URL.Query() {
		if len(value) > 0 {
			result = append(result, preservedParam{Key: key, Value: value[0]})
		}
	}
	return result
}
