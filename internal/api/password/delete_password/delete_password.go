package delete_password

import (
	"net/http"
	"os/exec"
	"log"
	"path/filepath"

	"github.com/gorilla/mux"

	show "pass_web/internal/api/show"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	passwordItem := show.PasswordsID[id]
	passwordFile := passwordItem.Password
	passwordPath := filepath.Join(passwordItem.Path, passwordFile)
	log.Println(passwordPath)

	cmd := exec.Command("rm", passwordPath)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to delete password", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
