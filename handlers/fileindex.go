package handlers

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/template"
	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/vfsutil"
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
	breadcrumb
	List []fileEntry
}

type fileindexHandler struct {
	root http.FileSystem
	path string
	endpoint config.FileindexHandlerEndpointStruct
	logger log.Channels
}
func FileindexHandler(
	root http.FileSystem,
	path string,
	endpoint config.FileindexHandlerEndpointStruct,
	logger log.Channels,
) http.Handler {
	template.AssertExists("fileindex", logger)
	return &fileindexHandler{root, path, endpoint, logger}
}

func (h *fileindexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tplData := fileindexPage {
		breadcrumb: prepareBreadcrum(r),
	}

	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	if !config.IsAllowedByACL(remoteAddr, h.endpoint.View) {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	if !assertPathBeginning(h.path, w, r, h.logger.Err) {
		return
	}

	if !config.Authenticate(r, h.endpoint.Auth) {
		authHeader := fmt.Sprintf(`Basic realm="Authentication required to use %s"`, tplData.LastBreadcrumb)
		w.Header().Set("WWW-Authenticate", authHeader)
		writeUnauthorizedError(w, r, h.logger.Err)
		return
	}

	if uploaded, err := h.uploadFile(w, r); uploaded {
		return
	} else if err != nil {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	if list, err := h.prepareFileList(r); err != nil {
		writeError(w, r, err, h.logger.Err)
		return
	} else {
		tplData.List = list
	}

	err := minifyTemplate("fileindex", tplData, w)
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

	if stat.Mode().IsDir() && req.URL.Query().Get("type") == "tar" {
		return h.uploadDir(writer, req)
	}

	if !stat.Mode().IsRegular() {
		return false, nil
	}

	http.ServeContent(writer, req, stat.Name(), stat.ModTime(), file)
	return true, nil
}

func (h *fileindexHandler) uploadDir(writer http.ResponseWriter, req *http.Request) (bool, error) {
	dir := req.URL.Path[1:]
	filename := fmt.Sprintf("%s.tar", filepath.Base(req.URL.Path))

	writer.Header().Add("Content-Type", "application/x-tar")
	writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	tw := tar.NewWriter(writer)
	defer tw.Close()

	fsroot := filter.Skip(h.root, func (path string, fi os.FileInfo) bool {
		return h.isHiddenPath(path)
	})

	err := vfsutil.Walk(fsroot, dir, func (path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		th, err := tar.FileInfoHeader(info, "") // TODO: better link argument?
		if err != nil {
			return err
		}
		th.Name = strings.TrimLeft(path[len(dir):], "/")
		if len(th.Name) == 0 { // base directory
			return err
		}

		fh, err := h.root.Open(path)
		if err != nil {
			return err
		}
		defer fh.Close()

		if err = tw.WriteHeader(th); err != nil {
			return err
		}
		if info.IsDir() {
			return err
		}

		_, err = io.Copy(tw, fh)
		return err
	})

	return err != nil, err
}

func (h *fileindexHandler) isHiddenPath(p string) bool {
	hidden := config.Get().Handlers.FileIndex.Hide
	dirname, filename := path.Split(p)
	for _, hiddenEntry := range hidden {
		hiddenFolder := fmt.Sprintf("/%s/", hiddenEntry)
		if filename == hiddenEntry || strings.Contains(dirname, hiddenFolder) {
			return true
		}
	}
	return false
}
