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
	conf := config.Get()
	basePath := conf.Paths.Templates
	if !filepath.IsAbs(basePath) {
		basePath = filepath.Join(conf.Paths.Base, basePath)
	}
	templates = template.New("")

	counter := 0
	root := filter.Keep(http.Dir(basePath), func (path string, fi fs.FileInfo) bool {
		return fi.IsDir() || filepath.Ext(path) == suffix
	})
	err := vfsutil.WalkFiles(root, "", func (path string, fi fs.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			log.Stderr().Print(err)
			return nil
		}

		if path == "" || fi.IsDir() {
			return nil
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			log.Stderr().Print(err)
			return nil
		}

		path = path[:len(path) - len(suffix)]
		_, err = templates.New(path).Parse(string(content))
		if err != nil {
			log.Stderr().Print(err)
		} else {
			counter += 1
		}

		return nil
	})

	if err != nil {
		log.Stderr().Fatal(err)
	}

	log.Stdout().Printf("Total templates at '%s': %v", basePath, counter)
}

func Execute(name string, data interface{}, w io.Writer) error {
	return templates.ExecuteTemplate(w, name, data)
}

func AssertExists(name string) {
	if !lo.ContainsBy(templates.Templates(), func (tpl *template.Template) bool { return tpl.Name() == name }) {
		log.Stderr().Fatalf("'%s' template is missing!", name)
	}
}
