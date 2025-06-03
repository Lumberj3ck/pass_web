package show

import (
	"html/template"
	"log"
	"net/http"
	"os/exec"
	templ "pass_web/internal/api/template"
	"path/filepath"
	"regexp"
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

		cmd := exec.Command("pass")
		output, err := cmd.Output()

		if err != nil {
			log.Println("cmd.Run() failed with %s\n", err)
		}
		var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
		output_str := ansiRegex.ReplaceAllString(string(output), "")

		lines := strings.Split(string(output_str), "\n")
		t.Execute(w, Page{lines})
	}
}
