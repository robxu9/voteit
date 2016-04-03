package routers

import (
	"html/template"
	"time"

	"github.com/unrolled/render"
)

var (
	renderer = render.New(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".tpl"},
		Funcs: []template.FuncMap{
			template.FuncMap{
				"rfc3339": func(t time.Time) string {
					return t.Format(time.RFC3339)
				},
			},
		},
		IndentJSON:    true,
		IndentXML:     true,
		IsDevelopment: true,
	})
)
