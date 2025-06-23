package insert_password

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"os"
	templ "pass_web/internal/api/template"
	"path/filepath"
)

func Handler(te *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost{
			password_name := r.FormValue("password_name")
			password_encrypted := r.FormValue("password_encrypted")


			if password_name == "" || password_encrypted == "" {
				http.Error(w, "Password name and password are required", http.StatusBadRequest)
				return
			}

			// if !strings.Contains(password_encrypted, "-----END PGP MESSAGE-----") {
			// 	http.Error(w, "Password is not encrypted", http.StatusBadRequest)
			// 	return
			// }
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
			// file.WriteString(password_encrypted)
			// file.WriteString(string(decodedBytes))
			file.Write(decodedBytes)

			te.Render(w, "password-insert-success", struct{}{})
		} else {
			dir := templ.GetTemplateDir()
			t, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir, "insert_password.tmpl"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			t.Execute(w, nil)
		}
	}
}
