package main

import (
	"embed"
	"log"
	"net/http"

	router "pass_web/internal/api/router"
	templ "pass_web/internal/api/template"
)

//go:embed templates/*
var TemplateFS embed.FS

//go:embed static/*
var clientAssets embed.FS

func main() {
	templ.TemplateFS = TemplateFS

	mu := router.NewMutexHandler()

	fs := http.FileServer(http.FS(clientAssets))
	mu.PathPrefix("/static/").Handler(fs)

	log.Println("Started listening on port 8080")
	err := http.ListenAndServe(":8080", mu)
	if err != nil {
		panic(err)
	}
}
