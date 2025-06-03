package show

import (
	"html/template"
	"log"
	"net/http"
	"os/exec"
	templ "pass_web/internal/api/template"
	"path/filepath"
	"strings"
)

type Page struct {
	Passwords []string
}

func Handler(t *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, e *http.Request) {
		dir := templ.GetTemplateDir()
		t, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir, "show.tmpl"))

		if err != nil {
			panic(err)
		}

		cmd := exec.Command("ls /root/.password-store")
		output, err := cmd.Output()

		if err != nil {
			log.Printf("cmd.Run() failed with %s\n", err)
		}


		lines := strings.Split(string(output), "\n")
		t.Execute(w, Page{lines})
	}
}
