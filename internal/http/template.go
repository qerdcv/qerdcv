package http

import (
	"html/template"
	"io"
)

type Template struct {
	template *template.Template
}

func NewTemplate(template *template.Template) *Template {
	return &Template{
		template: template,
	}
}

func (t *Template) Render(w io.Writer, name string, data any) error {
	return t.template.ExecuteTemplate(w, name, data)
}
