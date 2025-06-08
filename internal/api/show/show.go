package show

import (
	"html/template"
	"log"
	"net/http"
	"os/exec"
	templ "pass_web/internal/api/template"
	auth "pass_web/internal/api/auth"
	"path/filepath"
	"strings"
)


type PasswordItem struct{
	Id string
	Password string
}

var PasswordsID map[string]PasswordItem

type Page struct {
	Passwords []PasswordItem
}

func init(){
	PasswordsID = make(map[string]PasswordItem)
}

func Handler(t *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, e *http.Request) {
		dir := templ.GetTemplateDir()
		t, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir, "show.tmpl"))

		if err != nil {
			panic(err)
		}

		cmd := exec.Command("ls", "/root/.password-store")
		output, err := cmd.Output()

		if err != nil {
			log.Printf("cmd.Run() failed with %s\n", err)
		}

		lines := strings.Split(string(output), "\n")
		lines = lines[:len(lines) - 1]
		p := Page{}
		for i := 0; i < len(lines); i++{
			passwordID := auth.GenerateChallenge(20)
			p.Passwords = append(p.Passwords, PasswordItem{passwordID, lines[i]})
			PasswordsID[passwordID] = PasswordItem{passwordID, lines[i]}
		}
		t.Execute(w, p)
	}
}
