package router

import (
	"log"
	"net/http"
	"os"
	"pass_web/internal/api/auth"
	"pass_web/internal/api/show"
	templ "pass_web/internal/api/template"
	"github.com/gorilla/mux"
	password "pass_web/internal/api/password"
	"path/filepath"
)

func NewMutexHandler() *mux.Router {
	templ := templ.NewTemplate()
	log.Println("Created a new template handler")
	mu := mux.NewRouter()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mu.Handle("/", http.FileServer(http.Dir(filepath.Join(wd, "static"))))
	mu.HandleFunc("/show", auth.AuthMiddlerware(show.Handler(&templ)))
	mu.HandleFunc("/auth", auth.Handler(&templ))
	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(password.Handler(&templ)))

	log.Println("Started listening")
	return mu
}
