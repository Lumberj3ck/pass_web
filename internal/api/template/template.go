package api

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"embed"

	"github.com/joho/godotenv"
)

const PROJECT_ROOT_ENV = "PASS_WEB_ROOT"

var TemplateFS embed.FS

func GetTemplateDir() string {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dir := os.Getenv(PROJECT_ROOT_ENV)
	if dir == "" {
		fmt.Println("Error: PROJECT_ROOT environment variable is not set")
		os.Exit(1)
	}

	return filepath.Join(dir, "templates")
}

func NewTemplate(files ...string) Template {
	// dir := GetTemplateDir()
	var t *template.Template
	if len(files) > 0{
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
	if name == ""{
		t.tmpl.Execute(w, data)
		return 
	}
	t.tmpl.ExecuteTemplate(w, name, data)
}
