package show_password

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"os"

	show "pass_web/internal/api/show"

	"github.com/gorilla/mux"

	templ "pass_web/internal/api/template"
)

func Handler(w http.ResponseWriter, r *http.Request) {
		t := templ.NewTemplate()
		id := mux.Vars(r)["id"]
		
		passwordItem := show.PasswordsID[id]
		passwordFile := passwordItem.Password

		passwordPath := passwordItem.Path
		passwordPath = filepath.Join(passwordPath, passwordItem.Password)

		log.Println("Path password ", passwordPath)
		file, err := os.Open(passwordPath)

		if err != nil{
			http.Error(w, "Failed to show password; Password inaccesible", http.StatusInternalServerError)
		}
		password_buffer, _ := io.ReadAll(file)

		encodedContent := base64.StdEncoding.EncodeToString(password_buffer)

		w.Header().Set("Content-Type", "text/html")

		t.Render(w, "password", struct {
			PasswordFile           string
			EncodedContent string
		}{
			PasswordFile :           passwordFile,
			EncodedContent: encodedContent,
		})
	}
