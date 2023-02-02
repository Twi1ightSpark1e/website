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

type uploader func(w http.ResponseWriter, r *http.Request) (bool, error)

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
		"tar": func (w http.ResponseWriter, r *http.Request) (bool, error) { return h.uploadTar(w, r) },
		"gz": func (w http.ResponseWriter, r *http.Request) (bool, error) { return h.uploadGz(w, r) },
		"zst": func (w http.ResponseWriter, r *http.Request) (bool, error) { return h.uploadZst(w, r) },
	}

	return h
}

type page struct {
	util.BreadcrumbData
	markdown.InlineMarkdown
	findParams

	URL string
	AllowDownload bool
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
		AllowDownload: true,
		AllowUpload: allowUpload,
		URL: r.URL.Path,
		findParams: findParams{
			FindQuery: r.URL.Query().Get("query"),
			FindMatchCase: r.URL.Query().Get("matchcase") == "on",
			FindRegex: r.URL.Query().Get("regex") == "on",
		},
	}
	pageData.AllowDownload = pageData.AllowDownload && len(pageData.findParams.FindQuery) == 0
	pageData.AllowUpload = pageData.AllowUpload && len(pageData.findParams.FindQuery) == 0

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

	if sent, err := h.sendFile(w, r); sent {
		return
	} else if err != nil {
		errors.WriteNotFoundError(w, r, h.logger.Err)
		return
	}

	if list, err := h.prepareFileList(r.URL.Path, remoteAddr, pageData.findParams); err != nil {
		errors.WriteError(w, r, err, h.logger.Err)
		return
	} else {
		pageData.List = list

		show, name := h.showMarkdown(list)
		ptype := config.PreviewNone
		if show {
			ptype = h.endpoint.Preview
		}
		path := fmt.Sprintf("%s/%s", r.URL.Path, name)
		file, _ := h.root.Open(path)
		pageData.InlineMarkdown = markdown.PrepareInline(ptype, file)
		file.Close()
	}

	err := util.MinifyTemplate("fileindex", pageData, w)
	if err != nil {
		h.logger.Err.Print(err)
	}
}
