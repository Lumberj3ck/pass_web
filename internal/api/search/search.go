package search

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	show "pass_web/internal/api/show"
	templ "pass_web/internal/api/template"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	query := r.Form.Get("query")

	if query == "" {
		return
	}

	t := templ.NewTemplate("templates/password-item.tmpl")

	prefix := os.Getenv("PREFIX")

	fuzzy_entries := make([]string, 0)

	err = filepath.WalkDir(prefix, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel_path, err := filepath.Rel(prefix, path)
		if err != nil {
			return err
		}

		if strings.HasPrefix(rel_path, ".") {
			return nil
		}

		if !d.IsDir() {
			rel_path, err := filepath.Rel(prefix, path)
			if err != nil {
				return err
			}
			fuzzy_entries = append(fuzzy_entries, rel_path)
		}
		return nil
	})
	matches := fuzzy.Find(query, fuzzy_entries)

	var resp strings.Builder

	for _, match := range matches {
		var pi show.PasswordItem
		puid, ok := show.PasswordsPath[match]

		if !ok {
			pi = show.NewPasswordItem(filepath.Base(match), false, filepath.Dir(filepath.Join(prefix, match)), match)
			show.PasswordsID[pi.Id] = pi
			show.PasswordsPath[match] = pi.Id
		} else {
			pi = show.PasswordsID[puid]
		}

		pageItem := show.PasswordPageItem{
			PasswordItem: pi,
			Relative:     true,
			WithConsume:  false,
		}

		t.Render(&resp, "password-item", pageItem)
	}
	w.Write([]byte(resp.String()))
}
