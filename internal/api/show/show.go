package show

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	auth "pass_web/internal/api/auth"
	templ "pass_web/internal/api/template"
	utils "pass_web/internal/utils"
	"path/filepath"
	"strings"
	"sync"
)

type PasswordItem struct {
	Id           string
	Password     string
	IsDir        bool
	Path         string
	RelativePath string
}

func NewPasswordItem(password string, isDir bool, path string, relativePath string) PasswordItem {
	id := auth.GenerateChallenge(20)
	return PasswordItem{id, password, isDir, path, relativePath}
}

type PasswordIdStore struct {
	passwordsID   map[string]PasswordItem
	passwordsPath map[string]string
	mu            sync.RWMutex
}

func NewPasswordIdStore() *PasswordIdStore {
	return &PasswordIdStore{
		passwordsID:   make(map[string]PasswordItem),
		passwordsPath: make(map[string]string),
	}
}

func (p *PasswordIdStore) Add(password PasswordItem) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.passwordsID[password.Id] = password
	p.passwordsPath[password.Password] = password.Id
}

func (p *PasswordIdStore) Get(id string) (PasswordItem, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pi, ok := p.passwordsID[id]
	return pi, ok
}

func (p *PasswordIdStore) GetUid(path string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	id, ok := p.passwordsPath[path]
	return id, ok
}

func (p *PasswordIdStore) GetByPassword(password string) (PasswordItem, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pi, ok := p.passwordsPath[password]
	return p.passwordsID[pi], ok
}

type PasswordPageItem struct {
	PasswordItem
	Relative    bool
	WithConsume bool
}

type Page struct {
	Is_root   bool
	Passwords []PasswordPageItem
}

func Handler(ps *PasswordIdStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := templ.NewTemplate("templates/base.tmpl", "templates/show.tmpl", "templates/password-item.tmpl")

		prefix := utils.GetStorePrefix()
		uri_params := r.URL.Query()

		folder_id := uri_params["folder-id"]

		password_path := prefix
		is_root := true
		if len(folder_id) > 0 {
			pi, ok := ps.Get(folder_id[0])
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

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

			rel_path, err := filepath.Rel(prefix, filepath.Join(password_path, entry.Name()))

			if err != nil {
				slog.Warn("Failed to get relative path", "error", err)
				continue
			}
			pi, ok := ps.GetByPassword(rel_path)
			if !ok {
				pi = NewPasswordItem(entry.Name(), entry.IsDir(), password_path, rel_path)
				ps.Add(pi)
			}
			pageItem := PasswordPageItem{
				pi,
				false,
				true,
			}
			p.Passwords = append(p.Passwords, pageItem)
		}
		t.Render(w, "", p)
	}
}
