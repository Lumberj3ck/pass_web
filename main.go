package main

import (
	"net/http"
	"os"
	"log"
	"embed"
	"path/filepath"

	router "pass_web/internal/api/router"
	templ "pass_web/internal/api/template"

)

//go:embed templates/*
var TemplateFS embed.FS

func main() {
	templ.TemplateFS = TemplateFS
	wd := os.Getenv("PASS_WEB_ROOT")
	if wd == "" {
		wd = "."
	}

	mu := router.NewMutexHandler()
	staticDir := filepath.Join(wd, "static")
	fs := http.FileServer(http.Dir(staticDir))
	mu.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	
		
	log.Println("Started listening on port 8080")
	err := http.ListenAndServe(":8080", mu)
	if err != nil {
		panic(err)
	}
}
