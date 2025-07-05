package main

import (
	"net/http"
	"os"
	"log"
	"path/filepath"

	router "pass_web/internal/api/router"
)

func main() {
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
