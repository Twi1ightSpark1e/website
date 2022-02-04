package template

import (
	"fmt"
	"html/template"
	"io"

	"github.com/Twi1ightSpark1e/website/config"
	"github.com/Twi1ightSpark1e/website/log"
	"github.com/Twi1ightSpark1e/website/util"
)

var suffix = ".tpl"
var templates *template.Template

func templateFullName(name string) string {
	return fmt.Sprintf("%s%s", name, suffix)
}

func Initialize() {
	logger := log.New("TemplatesParser")

	basePath := util.FullPath(config.Get().TemplatesPath)
	logger.Info.Printf("Parsing templates at '%s'", basePath)

	var err error
	files, err := util.Glob(basePath, suffix)
	if err != nil {
		logger.Err.Fatal(err)
	}

	templates, err = template.ParseFiles(files...)
	if err != nil {
		logger.Err.Fatal(err)
	}
	for _, template := range templates.Templates() {
		fullName := template.Name()
		name := fullName[:len(fullName) - len(suffix)]
		logger.Info.Printf("New template registered: '%s'", name)
	}
}

func Execute(name string, data interface{}, w io.Writer) error {
	return templates.ExecuteTemplate(w, templateFullName(name), data)
}

func AssertExists(name string, logger log.Channels) {
	for _, template := range templates.Templates() {
		fullName := template.Name()
		templateName := fullName[:len(fullName) - len(suffix)]
		if templateName == name {
			return
		}
	}

	logger.Err.Fatalf("'%s' template is missing!", name)
}
