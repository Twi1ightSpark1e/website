package template

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
)

var templates map[string]*template.Template

func Initialize() {
	templates = make(map[string]*template.Template)

	basePath := "template/"
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
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
	}
}

func Get(name string) *template.Template {
	return templates[name]
}
