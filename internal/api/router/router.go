package router

import (
	"log"
	"net/http"
	"os"
	"pass_web/internal/api/auth"
	"pass_web/internal/api/show"
	templ "pass_web/internal/api/template"
	"path/filepath"
)


func NewMutexHandler() *http.ServeMux{
    templ := templ.NewTemplate()
    log.Println("Created a new template")
    mu := http.NewServeMux()

    wd, err := os.Getwd()
    if err != nil{
        panic(err)
    }

    mu.Handle("/", http.FileServer(http.Dir(filepath.Join(wd, "static"))))
    mu.HandleFunc("/show", show.Handler(&templ))
    mu.HandleFunc("/auth", auth.Handler(&templ))

    log.Println("Started listening")
    return mu
}
