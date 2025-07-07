package router

import (
	"log"
	"pass_web/internal/api/auth"
	"pass_web/internal/api/show"

	delete_password "pass_web/internal/api/password/delete_password"
	insert_password "pass_web/internal/api/password/insert_password"
	show_password "pass_web/internal/api/password/show_password"

	"github.com/gorilla/mux"
)

func NewMutexHandler() *mux.Router {
	log.Println("Created a new template handler")
	mu := mux.NewRouter()

	mu.HandleFunc("/show", auth.AuthMiddlerware(show.Handler))
	mu.HandleFunc("/auth", auth.Handler)

	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(show_password.Handler)).Methods("POST")
	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(delete_password.Handler)).Methods("DELETE")
	mu.HandleFunc("/insert", auth.AuthMiddlerware(insert_password.Handler)).Methods("POST", "GET")

	log.Println("Started listening")
	return mu
}
