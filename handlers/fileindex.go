package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

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

type breadcrumbItem struct {
	Title string
	Address string
}

type fileindexPage struct {
	Breadcrumb []breadcrumbItem
	LastBreadcrumb string
	List []fileEntry
	Error string
}

func FileindexHandler(w http.ResponseWriter, r *http.Request) {
	breadcrumb := prepareBreadcrum(r)
	tplData := fileindexPage {
		Breadcrumb: breadcrumb[:len(breadcrumb) - 1],
		LastBreadcrumb: breadcrumb[len(breadcrumb) - 1].Title,
	}

	upload, err := shouldUploadFile(r)
	if err != nil {
		tplData.Error = err.Error()
	} else if upload {
		err = uploadFile(w, r, tplData.LastBreadcrumb)
		if err == nil {
			return
		}
		tplData.Error = err.Error()
	} else {
		list, err := prepareFileList(r)
		if err != nil {
			tplData.Error = err.Error()
		} else {
			tplData.List = list
		}
	}

	err = template.Get("fileindex").Execute(w, tplData)
	if err != nil {
		log.Print(err)
	}
}

func prepareBreadcrum(req *http.Request) []breadcrumbItem {
	result := []breadcrumbItem {
		{
			Title: req.Host,
			Address: "/",
		},
	}

	items := filterStr(strings.Split(req.URL.Path, "/"), func (item *string) bool {
		return len(*item) != 0
	})
	for idx, item := range items {
		if len(item) == 0 {
			continue
		}

		address := fmt.Sprintf("/%s", strings.Join(items[:idx + 1], "/"))
		result = append(result, breadcrumbItem {
			Title: item,
			Address: address,
		})
	}

	return result
}

func prepareFileList(req *http.Request) ([]fileEntry, error) {
	result := make([]fileEntry, 0)

	_, filename := path.Split(req.URL.Path)
	if filename == "noindex" {
		return result, errors.New("You are not allowed to view this folder")
	}

	files, err := ioutil.ReadDir(fmt.Sprintf("..%c%s", os.PathSeparator, req.URL.Path))
	if err != nil {
		return result, err
	}

	for _, file := range files {
		name := file.Name()
		if name == "noindex" {
			continue
		}
		if file.IsDir() {
			name = name + string(os.PathSeparator)
		}

		entry := fileEntry {
			IsDir: file.IsDir(),
			Name: name,
			Date: file.ModTime().UTC().Format("2006-01-02 15:04:05"),
		}
		if !entry.IsDir {
			entry.Size = ByteCountIEC(file.Size())
		}

		result = append(result, entry)
	}

	err = nil
	if len(result) == 0 {
		err = errors.New("This folder is empty")
	}
	return result, err
}

func shouldUploadFile(req *http.Request) (bool, error) {
	stat, err := os.Stat(fmt.Sprintf("..%c%s", os.PathSeparator, req.URL.Path))
	if err != nil {
		return false, err
	}

	return stat.Mode().IsRegular(), nil
}

func uploadFile(w http.ResponseWriter, r *http.Request, filename string) error {
	path := fmt.Sprintf("..%c%s", os.PathSeparator, r.URL.Path)
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	stat, err := reader.Stat()
	if err != nil {
		return err
	}

	http.ServeContent(w, r, filename, stat.ModTime(), reader)
	return nil
}
