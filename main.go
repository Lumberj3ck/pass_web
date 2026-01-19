package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	router "pass_web/internal/api/router"
	show "pass_web/internal/api/show"
	templ "pass_web/internal/api/template"
)

//go:embed templates/*
var TemplateFS embed.FS

//go:embed static/*
var clientAssets embed.FS

func main() {
	templ.TemplateFS = TemplateFS
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	passwordStore := show.NewPasswordIdStore()
	mu := router.NewMutexHandler(passwordStore)

	fs := http.FileServer(http.FS(clientAssets))
	mu.PathPrefix("/static/").Handler(fs)

	log.Println("Started http server ", "port", *port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), mu)
	if err != nil {
		panic(err)
	}
}
