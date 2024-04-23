package web

import (
	"algo/models"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Playground struct {
	Editor         models.Editor
	AllPlaygrounds []string
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func GetRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob("web/*.html")),
	}
}
