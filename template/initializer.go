package template

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/util"
)

var templates map[string]*template.Template

func Initialize() {
	templates = make(map[string]*template.Template)
	logger := log.New("TemplatesParser")

	basePath := util.FullPath(config.Get().TemplatesPath)
	logger.Info.Printf("Parsing templates at '%s'", basePath)
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		logger.Err.Fatal(err)
		return
	}

	suffix := ".tpl"
	for _, f := range files {
		name := f.Name()
		if !strings.HasSuffix(name, suffix) {
			continue
		}

		path := fmt.Sprintf("%s%s", basePath, name)
		tpl := template.Must(template.ParseFiles(path))
		templateName := name[:len(name) - len(suffix)]
		templates[templateName] = tpl

		logger.Info.Printf("New template registered: '%s'", templateName)
	}
}

func Get(name string) *template.Template {
	return templates[name]
}

func AssertExists(name string, logger log.Channels) {
	if templates[name] == nil {
		logger.Err.Fatalf("'%s' template is missing!", name)
	}
}
