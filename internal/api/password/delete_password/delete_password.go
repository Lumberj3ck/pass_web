package delete_password

import (
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"

	show "pass_web/internal/api/show"

	templ "pass_web/internal/api/template"
)

func Handler(t *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		passwordItem := show.PasswordsID[id]
		passwordFile := passwordItem.Password
		passwordPath := filepath.Join("/root/.password-store", passwordFile)

		cmd := exec.Command("rm", passwordPath)
		err := cmd.Run()
		if err != nil {
			http.Error(w, "Failed to delete password", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}