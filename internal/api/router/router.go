package router

import (
	"log"
	"pass_web/internal/api/auth"
	"pass_web/internal/api/show"
	templ "pass_web/internal/api/template"

	password "pass_web/internal/api/password"

	"github.com/gorilla/mux"
)

func NewMutexHandler() *mux.Router {
	templ := templ.NewTemplate()
	log.Println("Created a new template handler")
	mu := mux.NewRouter()


	mu.HandleFunc("/show", auth.AuthMiddlerware(show.Handler(&templ)))
	mu.HandleFunc("/auth", auth.Handler(&templ))
	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(password.Handler(&templ)))

	log.Println("Started listening")
	return mu
}
