package api

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const PROJECT_ROOT_ENV = "PASS_WEB_ROOT"

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

	fmt.Println(dir)
	return filepath.Join(dir, "templates")
}

func NewTemplate() Template {
	dir := GetTemplateDir()
	return Template{
		tmpl: template.Must(template.ParseGlob(filepath.Join(dir, "*.tmpl"))),
	}
}

type Template struct {
	tmpl *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}) {
	t.tmpl.ExecuteTemplate(w, name, data)
}
