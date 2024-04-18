package server

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"

	"github.com/qerdcv/qerdcv/web"
)

type TemplateRenderer struct {
}

func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{}
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	if viewContext, ok := data.(map[string]any); ok {
		viewContext["reverse"] = c.Echo().Reverse
	}

	baseTmpl := template.Must(template.ParseFS(web.Templates, "templates/base.gohtml"))
	baseTmpl = template.Must(baseTmpl.ParseFS(web.Templates, name))
	return baseTmpl.Execute(w, data)
}
