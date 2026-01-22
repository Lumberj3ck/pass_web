package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	"pass_web/internal/api/auth"
	router "pass_web/internal/api/router"
	render_folder "pass_web/internal/api/render_folder"
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

	passwordStore := render_folder.NewPasswordIdStore()
	userChallenges := auth.NewUserChalenges()
	mu := router.NewMutexHandler(passwordStore, userChallenges)

	fs := http.FileServer(http.FS(clientAssets))
	mu.PathPrefix("/static/").Handler(fs)

	log.Println("Started http server ", "port", *port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), mu)
	if err != nil {
		panic(err)
	}
}
