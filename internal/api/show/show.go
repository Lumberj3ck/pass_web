package show

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	auth "pass_web/internal/api/auth"
	templ "pass_web/internal/api/template"
	"path/filepath"
	"strings"
)

type PasswordItem struct {
	Id       string
	Password string
	IsDir    bool
	Path     string
}

var PasswordsID map[string]PasswordItem

type Page struct {
	Is_root   bool
	Passwords []PasswordItem
}

func init() {
	PasswordsID = make(map[string]PasswordItem)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	t := templ.NewTemplate("templates/base.tmpl", "templates/show.tmpl")

	prefix := os.Getenv("PREFIX")
	uri_params := r.URL.Query()

	folder_id := uri_params["folder-id"]

	password_path := prefix
	is_root := true
	if len(folder_id) > 0 {
		pi, ok := PasswordsID[folder_id[0]]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// password_path = filepath.Join(prefix, folder_id[0])
		password_path = filepath.Join(pi.Path, pi.Password)
		is_root = false
	}

	slog.Info("Show password ", "password_path", password_path)

	entries, err := os.ReadDir(password_path)

	if err != nil {
		log.Printf("Failed to read from password store: %s\n", err)
	}

	p := Page{}
	p.Is_root = is_root

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		passwordID := auth.GenerateChallenge(20)

		p.Passwords = append(p.Passwords, PasswordItem{passwordID, entry.Name(), entry.IsDir(), password_path})

		PasswordsID[passwordID] = PasswordItem{passwordID, entry.Name(), entry.IsDir(), password_path}
	}
	t.Render(w, "", p)
}
