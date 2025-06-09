package show_password

import (
	"encoding/base64"
	"log"
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

		cmd := exec.Command("cat", passwordPath)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Failed to list passwords: %v", err)
			http.Error(w, "Failed to list passwords", http.StatusInternalServerError)
			return
		}

		encodedContent := base64.StdEncoding.EncodeToString(output)

		w.Header().Set("Content-Type", "text/html")

		t.Render(w, "password", struct {
			PasswordFile           string
			EncodedContent string
		}{
			PasswordFile :           passwordFile,
			EncodedContent: encodedContent,
		})
	}
}