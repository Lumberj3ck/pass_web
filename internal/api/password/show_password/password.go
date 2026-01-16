package show_password

import (
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"

	"os"

	show "pass_web/internal/api/show"

	"github.com/gorilla/mux"

	templ "pass_web/internal/api/template"
)

type PasswordTempl struct {
	PasswordFile   string
	EncodedContent string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	t := templ.NewTemplate()
	id := mux.Vars(r)["id"]

	passwordItem := show.PasswordsID[id]
	passwordFile := passwordItem.Password

	passwordPath := passwordItem.Path
	passwordPath = filepath.Join(passwordPath, passwordFile)

	file, err := os.Open(passwordPath)

	if err != nil {
		http.Error(w, "Failed to show password; Password inaccesible", http.StatusInternalServerError)
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
