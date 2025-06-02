package show

import (
	"html/template"
	"net/http"
	templ "pass_web/internal/api/template"
	"path/filepath"
)

type Page struct{
    Password string
}

func Handler(t *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, e *http.Request) {
    
        dir := templ.GetTemplateDir()
        t, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir,"show.tmpl"))

        if err != nil{
            panic(err)
        }

        t.Execute(w, Page{"Passowrd"})


	}
}
