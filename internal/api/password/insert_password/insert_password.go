package insert_password

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	templ "pass_web/internal/api/template"
	"path/filepath"
)

func Handler(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost{
			t := templ.NewTemplate()
			password_name := r.FormValue("password_name")
			password_encrypted := r.FormValue("password_encrypted")


			if password_name == "" || password_encrypted == "" {
				http.Error(w, "Password name and password are required", http.StatusBadRequest)
				return
			}

			decodedBytes, err := base64.StdEncoding.DecodeString(password_encrypted)

			if err != nil {
				http.Error(w, "Failed to decode encrypted password (invalid Base64): "+err.Error(), http.StatusBadRequest)
				return
			}


			prefix := os.Getenv("PREFIX")
			log.Println(password_name)
			file, err := os.Create(filepath.Join(prefix, password_name))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()
			file.Write(decodedBytes)

			t.Render(w, "password-insert-success", struct{}{})
		} else {
			t := templ.NewTemplate("templates/base.tmpl", "templates/insert_password.tmpl")
			t.Render(w, "", nil)
		}
	}
