package router

import (
	"log"
	"pass_web/internal/api/auth"
	"pass_web/internal/api/render_folder"

	"pass_web/internal/api/password/delete_password"
	"pass_web/internal/api/password/insert_password"
	"pass_web/internal/api/password/show_password"
	"pass_web/internal/api/search"

	"github.com/gorilla/mux"
)

func NewMutexHandler(ps *render_folder.PasswordIdStore, uc *auth.UserChalenges) *mux.Router {
	log.Println("Created a new template handler")
	mu := mux.NewRouter()

	mu.HandleFunc("/", auth.AuthMiddlerware(render_folder.Handler(ps)))
	mu.HandleFunc("/auth", auth.Handler(uc))

	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(show_password.Handler(ps))).Methods("POST")
	mu.HandleFunc("/password/{id}", auth.AuthMiddlerware(delete_password.Handler(ps))).Methods("DELETE")
	mu.HandleFunc("/insert", auth.AuthMiddlerware(insert_password.Handler)).Methods("POST", "GET")

	mu.HandleFunc("/search", auth.AuthMiddlerware(search.Handler(ps)))
	log.Println("Started listening")
	return mu
}
