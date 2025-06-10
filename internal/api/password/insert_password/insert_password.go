package insert_password

import (
	"html/template"
	"log"
	"net/http"
	"os"
	templ "pass_web/internal/api/template"
	"path/filepath"
	"strings"
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

			if !strings.Contains(password_encrypted, "-----END PGP MESSAGE-----") {
				http.Error(w, "Password is not encrypted", http.StatusBadRequest)
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
			file.WriteString(password_encrypted)
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
