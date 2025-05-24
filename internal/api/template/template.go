package api

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

func GetTemplateDir() string {
    dir, err := os.Getwd()
    if err != nil{
        panic(err)
    }
    fmt.Println(dir)
    return filepath.Join(dir, "templates")
}

func NewTemplate() Template{
    dir := GetTemplateDir()
    return Template{
        tmpl: template.Must(template.ParseGlob(filepath.Join(dir, "*.tmpl"))),
    }
}

type Template struct{
    tmpl *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}) {
    t.tmpl.ExecuteTemplate(w, name, data)
}
