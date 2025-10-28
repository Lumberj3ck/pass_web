package api

import (
	"embed"
	"html/template"
	"io"
)

var TemplateFS embed.FS

func NewTemplate(files ...string) Template {
	var t *template.Template
	if len(files) > 0 {
		t = template.Must(template.ParseFS(TemplateFS, files...))
	} else {
		t = template.Must(template.ParseFS(TemplateFS, "templates/*.tmpl"))
	}
	return Template{
		tmpl: t,
	}
}

type Template struct {
	tmpl *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}) {
	if name == "" {
		t.tmpl.Execute(w, data)
		return
	}
	t.tmpl.ExecuteTemplate(w, name, data)
}
