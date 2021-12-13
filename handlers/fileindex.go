package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
)

type filterPred func (item *string) bool
func filterStr(data []string, predicate filterPred) []string {
	result := make([]string, 0)

	for _, item := range data {
		if predicate(&item) {
			result = append(result, item)
		}
	}

	return result
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

type fileEntry struct {
	Name string
	Size string
	Date string
	IsDir bool
}

type fileindexPage struct {
	Title string
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
	List []fileEntry
	Error string
}

type fileindexHandler struct {
	root http.FileSystem
	endpoint config.FileindexHandlerEndpointStruct
	logger log.Channels
}
func FileindexHandler(
	root http.FileSystem,
	endpoint config.FileindexHandlerEndpointStruct,
	logger log.Channels,
) http.Handler {
	return &fileindexHandler{root, endpoint, logger}
}

func (h *fileindexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	breadcrumb := prepareBreadcrum(r)
	tplData := fileindexPage {
		Title: prepareTitle(r),
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
	}

	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		w.WriteHeader(http.StatusNotFound)
		tplData.Error = "Content not found"
		goto compile
	}

	if uploaded, err := h.uploadFile(w, r); uploaded {
		return
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
		tplData.Error = "Content not found"
		goto compile
	}

	if list, err := h.prepareFileList(r); err != nil {
		w.WriteHeader(http.StatusNotFound)
		tplData.Error = err.Error()
	} else {
		tplData.List = list
	}

compile:
	err := template.Get("fileindex").Execute(w, tplData)
	if err != nil {
		h.logger.Err.Print(err)
	}
}

func (h *fileindexHandler) prepareFileList(req *http.Request) ([]fileEntry, error) {
	result := make([]fileEntry, 0)

	if h.isHiddenPath(req.URL.Path) {
		return result, errors.New("Content not found")
	}

	direntry, err := h.root.Open(req.URL.Path)
	if err != nil {
		return result, err
	}

	files, err := direntry.Readdir(-1)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		name := file.Name()
		if h.isHiddenPath(name) {
			continue
		}
		if file.IsDir() {
			name = name + string(os.PathSeparator)
		}

		result = append(result, fileEntry {
			IsDir: file.IsDir(),
			Name: name,
			Date: file.ModTime().UTC().Format("2006-01-02 15:04:05"),
			Size: ByteCountIEC(file.Size()),
		})
	}

	err = nil
	if len(result) == 0 {
		err = errors.New("This folder is empty")
	} else {
		sort.Slice(result, func (i, j int) bool {
			if result[i].IsDir != result[j].IsDir {
				return result[i].IsDir
			}
			return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
		})
	}
	return result, err
}

func (h *fileindexHandler) uploadFile(writer http.ResponseWriter, req *http.Request) (bool, error) {
	file, err := h.root.Open(req.URL.Path)
	if err != nil {
		return false, err
	}

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	if !stat.Mode().IsRegular() {
		return false, nil
	}

	http.ServeContent(writer, req, stat.Name(), stat.ModTime(), file)
	return true, nil
}

func (h *fileindexHandler) isHiddenPath(p string) bool {
	hidden := config.Get().Handlers.FileIndex.Hide
	dirname, filename := path.Split(p)
	for _, hiddenEntry := range hidden {
		if filename == hiddenEntry || strings.Contains(dirname, hiddenEntry) {
			return true
		}
	}
	return false
}

func prepareTitle(req *http.Request) string {
	breadcrumb := prepareBreadcrum(req)
	return fmt.Sprintf("%s %s", req.Host, breadcrumb[1].Title)
}
