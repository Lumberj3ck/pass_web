package show_password

import (
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"

	"os"

	"pass_web/internal/api/render_folder"

	"github.com/gorilla/mux"

	templ "pass_web/internal/api/template"
)

type PasswordTempl struct {
	PasswordFile   string
	EncodedContent string
}

func Handler(ps *render_folder.PasswordIdStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := templ.NewTemplate()
		id := mux.Vars(r)["id"]

		passwordItem, ok := ps.Get(id)
		if !ok {
			http.Error(w, "Password inaccesible", http.StatusBadRequest)
			return
		}
		passwordFile := passwordItem.Password

		passwordPath := passwordItem.Path
		passwordPath = filepath.Join(passwordPath, passwordFile)

		file, err := os.Open(passwordPath)

		if err != nil {
			slog.Info("Ivalid password id", "id", id)
			http.Error(w, "Failed to show password; Password inaccesible", http.StatusBadRequest)
			return
		}
		password_buffer, _ := io.ReadAll(file)

		encodedContent := base64.StdEncoding.EncodeToString(password_buffer)

		w.Header().Set("Content-Type", "text/html")

		baseDir := filepath.Base(passwordItem.Path)
		relativeFilename := filepath.Join(baseDir, passwordFile)
		slog.Info("Password show ", "baseDir", baseDir, "relativeFilename", relativeFilename)

		t.Render(w, "password", PasswordTempl{
			PasswordFile:   relativeFilename,
			EncodedContent: encodedContent,
		})
	}
}
