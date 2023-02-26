package template

import (
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/samber/lo"
	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/vfsutil"
)

var suffix = ".tpl"
var templates *template.Template

func Initialize() {
	logger := log.New("TemplatesParser")

	conf := config.Get()
	basePath := filepath.Join(conf.Paths.Base, conf.Paths.Templates)
	logger.Info.Printf("Parsing templates at '%s'", basePath)
	templates = template.New("")

	root := filter.Keep(http.Dir(basePath), func (path string, fi fs.FileInfo) bool {
		return fi.IsDir() || filepath.Ext(path) == suffix
	})
	err := vfsutil.WalkFiles(root, "", func (path string, fi fs.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			logger.Err.Print(err)
			return nil
		}

		if path == "" || fi.IsDir() {
			return nil
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			logger.Err.Print(err)
			return nil
		}

		path = path[:len(path) - len(suffix)]
		_, err = templates.New(path).Parse(string(content))
		if err != nil {
			logger.Err.Print(err)
		} else {
			logger.Info.Printf("New template registered: '%s'", path)
		}

		return nil
	})

	if err != nil {
		logger.Err.Fatal(err)
	}
}

func Execute(name string, data interface{}, w io.Writer) error {
	return templates.ExecuteTemplate(w, name, data)
}

func AssertExists(name string, logger log.Channels) {
	if !lo.ContainsBy(templates.Templates(), func (tpl *template.Template) bool { return tpl.Name() == name }) {
		logger.Err.Fatalf("'%s' template is missing!", name)
	}
}
