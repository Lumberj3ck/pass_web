package router

import (
	"log"
	"pass_web/internal/api/auth"
	"pass_web/internal/api/show"
	templ "pass_web/internal/api/template"

	delete_password "pass_web/internal/api/password/delete_password"
	insert_password "pass_web/internal/api/password/insert_password"
	show_password "pass_web/internal/api/password/show_password"

	"github.com/gorilla/mux"
)

func NewMutexHandler() *mux.Router {
	templ := templ.NewTemplate()
	log.Println("Created a new template handler")
	mu := mux.NewRouter()


	mu.HandleFunc("/show", auth.AuthMiddlerware(show.Handler(&templ)))
	mu.HandleFunc("/auth", auth.Handler(&templ))

	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(show_password.Handler(&templ))).Methods("POST")
	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(delete_password.Handler(&templ))).Methods("DELETE")
	mu.HandleFunc("/insert", auth.AuthMiddlerware(insert_password.Handler(&templ))).Methods("POST", "GET")

	log.Println("Started listening")
	return mu
}
