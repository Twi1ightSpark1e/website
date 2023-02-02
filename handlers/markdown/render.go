package markdown

import (
	"html/template"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func (h *handler) render() (template.HTML, error) {
	file, err := h.root.Open(h.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	md, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return Render(md), nil
}

func Render(content []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.Attributes
	parser := parser.NewWithExtensions(extensions)
	html := template.HTML(markdown.ToHTML(content, parser, nil))

	return html
}
