package delete_password

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"pass_web/internal/api/render_folder"
)

func Handler(ps *render_folder.PasswordIdStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		passwordItem, ok := ps.Get(id)
		if !ok {
			http.Error(w, "Password inaccesible", http.StatusBadRequest)
			return
		}
		passwordFile := passwordItem.Password
		passwordPath := filepath.Join(passwordItem.Path, passwordFile)
		log.Println(passwordPath)

		err := os.Remove(passwordPath)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Failed to delete password", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
