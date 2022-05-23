package handlers

import (
	"archive/tar"
	"errors"
	"fmt"
	"html/template"
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
	tpl "github.com/Twi1ightSpark1e/website/template"
	"github.com/Twi1ightSpark1e/website/util"
	"github.com/flytam/filenamify"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/vfsutil"
)

func useAsPreview(name string) bool {
	for _, entry := range config.Get().Handlers.FileIndex.Preview {
		if name == entry {
			return true
		}
	}
	return false
}

type uploader func(w http.ResponseWriter, r *http.Request) (bool, error)

type fileEntry struct {
	Name string
	Size string
	Date string
	IsDir bool
}

type fileindexPage struct {
	breadcrumb
	AllowUpload bool
	List []fileEntry

	ShowMarkdown bool
	MarkdownOnTop bool
	MarkdownTitle string
	MarkdownContent template.HTML
}

type fileindexHandler struct {
	root http.FileSystem
	path string
	endpoint config.FileindexHandlerEndpointStruct
	logger log.Channels
	uploaders map[string]uploader
}
func FileindexHandler(
	root http.FileSystem,
	path string,
	endpoint config.FileindexHandlerEndpointStruct,
	logger log.Channels,
) http.Handler {
	tpl.AssertExists("fileindex", logger)

	h := &fileindexHandler{root, path, endpoint, logger, map[string]uploader{}}
	h.uploaders = map[string]uploader {
		"tar": func (w http.ResponseWriter, r *http.Request) (bool, error) { return h.uploadTar(w, r) },
		"gz": func (w http.ResponseWriter, r *http.Request) (bool, error) { return h.uploadGz(w, r) },
		"zst": func (w http.ResponseWriter, r *http.Request) (bool, error) { return h.uploadZst(w, r) },
	}

	return h
}

func (h *fileindexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := getRemoteAddr(r)
	h.logger.Info.Printf("Client %s requested '%s'", remoteAddr, r.URL.Path)

	allowUpload := config.IsAllowedByACL(remoteAddr, h.endpoint.Upload)
	allowPost := r.Method == http.MethodPost && allowUpload
	allowView := r.Method != http.MethodPost && config.IsAllowedByACL(remoteAddr, h.endpoint.View)
	if !allowPost && !allowView {
		writeNotFoundError(w, r, h.logger.Err)
		return
	}

	if !assertPathBeginning(h.path, w, r, h.logger.Err) {
		return
	}

	tplData := fileindexPage {
		breadcrumb: prepareBreadcrum(r),
		AllowUpload: allowUpload,
	}
	if !config.Authenticate(r, h.endpoint.Auth) {
		authHeader := fmt.Sprintf(`Basic realm="Authentication required to use %s"`, tplData.LastBreadcrumb)
		w.Header().Set("WWW-Authenticate", authHeader)
		writeUnauthorizedError(w, r, h.logger.Err)
		return
	}

	if recv, err := h.recvFile(w, r); recv {
		return
	} else if err != nil {
		writeError(w, r, err, h.logger.Err)
		return
	}

	if sent, err := h.sendFile(w, r); sent {
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

		show, name := h.showMarkdown(list)
		tplData.ShowMarkdown = show
		tplData.MarkdownOnTop = h.endpoint.Preview == config.PreviewPre
		tplData.MarkdownTitle = name
		if show {
			path := fmt.Sprintf("%s/%s", r.URL.Path, name)
			md, err := h.loadMarkdown(path)
			if err != nil {
				tplData.ShowMarkdown = false
			} else {
				tplData.MarkdownContent = md
			}
		}
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
			Size: util.ByteCountIEC(file.Size()),
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

func (h *fileindexHandler) sendFile(writer http.ResponseWriter, req *http.Request) (bool, error) {
	file, err := h.root.Open(req.URL.Path)
	if err != nil {
		return false, err
	}

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	if stat.Mode().IsDir() {
		if uploader, ok := h.uploaders[req.URL.Query().Get("type")]; ok {
			return uploader(writer, req)
		}
	}

	if !stat.Mode().IsRegular() {
		return false, nil
	}

	http.ServeContent(writer, req, stat.Name(), stat.ModTime(), file)
	return true, nil
}

func (h *fileindexHandler) prepareTar(w io.WriteCloser, dir string) error {
	defer w.Close()

	tw := tar.NewWriter(w)
	defer tw.Close()

	fsroot := filter.Skip(h.root, func (path string, fi os.FileInfo) bool {
		return h.isHiddenPath(path)
	})

	return vfsutil.Walk(fsroot, dir, func (path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		th, err := tar.FileInfoHeader(info, "") // TODO: better link argument?
		if err != nil {
			return err
		}

		// hide uid:gid, set them to nobody
		th.Uid = 65534
		th.Gid = 65534
		th.Uname = "nobody"
		th.Gname = "nobody"

		// fix file path
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
}

func (h *fileindexHandler) uploadTar(w http.ResponseWriter, r *http.Request) (bool, error) {
	dir := r.URL.Path[1:]
	filename := fmt.Sprintf("%s.tar", filepath.Base(r.URL.Path))

	w.Header().Add("Content-Type", "application/x-tar")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	bufr, bufw := io.Pipe()
	defer bufr.Close()

	go h.prepareTar(bufw, dir)
	written, err := io.Copy(w, bufr)
	return written > 0, err
}

func (h *fileindexHandler) uploadGz(w http.ResponseWriter, r *http.Request) (bool, error) {
	dir := r.URL.Path[1:]
	filename := fmt.Sprintf("%s.tar.gz", filepath.Base(r.URL.Path))

	w.Header().Add("Content-Type", "application/gzip")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	bufr, bufw := io.Pipe()
	defer bufr.Close()

	compressor := gzip.NewWriter(w)
	defer compressor.Close()

	go h.prepareTar(bufw, dir)
	written, err := io.Copy(compressor, bufr)
	return written > 0, err
}

func (h *fileindexHandler) uploadZst(w http.ResponseWriter, r *http.Request) (bool, error) {
	dir := r.URL.Path[1:]
	filename := fmt.Sprintf("%s.tar.zst", filepath.Base(r.URL.Path))

	w.Header().Add("Content-Type", "application/zstd")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	bufr, bufw := io.Pipe()
	defer bufr.Close()

	compressor, err := zstd.NewWriter(w)
	if err != nil {
		return false, err
	}
	defer compressor.Close()

	go h.prepareTar(bufw, dir)
	written, err := io.Copy(compressor, bufr)
	return written > 0, err
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

func (h *fileindexHandler) recvFile(w http.ResponseWriter, r * http.Request) (bool, error) {
	if r.Method != http.MethodPost {
		return false, nil
	}

	r.ParseMultipartForm(1024 * 1024)

	file, header, err := r.FormFile("file")
	if err != nil {
		return false, errors.New("No file chosen")
	}
	defer file.Close()

	filename, err := filenamify.Filenamify(header.Filename, filenamify.Options{
		Replacement: "_",
	})
	if err != nil {
		return false, err
	}
	filepath := r.URL.Path + filename
	h.logger.Info.Printf("Receiving file '%s'", filepath)

	destFile, err := os.Create(config.Get().Paths.Base + filepath)
	if err != nil {
		return false, err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		return false, err
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	return true, nil
}

func (h *fileindexHandler) showMarkdown(list []fileEntry) (bool, string) {
	for _, file := range list {
		if file.IsDir || !useAsPreview(file.Name) {
			continue
		}
		return true, file.Name
	}
	return false, ""
}

func (h *fileindexHandler) loadMarkdown(path string) (template.HTML, error) {
	file, err := h.root.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return renderMarkdown(buf), nil
}
