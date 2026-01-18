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

func NewPasswordItem(password string, isDir bool, path string) PasswordItem {
	id := auth.GenerateChallenge(20)
	return PasswordItem{id, password, isDir, path}
}

var PasswordsID map[string]PasswordItem
var PasswordsPath map[string]string

type Page struct {
	Is_root   bool
	Passwords []PasswordItem
}

func init() {
	PasswordsID = make(map[string]PasswordItem)
	PasswordsPath = make(map[string]string)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	t := templ.NewTemplate("templates/base.tmpl", "templates/show.tmpl", "templates/password-item.tmpl")

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

		var pi PasswordItem

		if _, ok := PasswordsPath[entry.Name()]; !ok{
			pi = NewPasswordItem(entry.Name(), entry.IsDir(), password_path)
			PasswordsID[pi.Id] = pi 
			PasswordsPath[entry.Name()] = pi.Id
		} else {
			pi = PasswordsID[PasswordsPath[entry.Name()]]
		}
		p.Passwords = append(p.Passwords, pi)
	}
	t.Render(w, "", p)
}
