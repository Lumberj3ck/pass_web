package search

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
		fmt.Fprintf(&resp, "<div>%v</div>", match)
	}
	w.Write([]byte(resp.String()))
}
