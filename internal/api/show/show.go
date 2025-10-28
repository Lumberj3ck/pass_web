package show

import (
	"log"
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

func Handler(w http.ResponseWriter, e *http.Request) {
	t := templ.NewTemplate("templates/base.tmpl", "templates/show.tmpl")

	prefix := os.Getenv("PREFIX")
	uri_params := e.URL.Query()

	subpath := uri_params["path"]

	password_path := prefix
	is_root := true
	if len(subpath) > 0 {
		password_path = filepath.Join(prefix, subpath[0])
		is_root = false
	}
	log.Println(password_path)

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
