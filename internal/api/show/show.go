package show

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	auth "pass_web/internal/api/auth"
	templ "pass_web/internal/api/template"
	"path/filepath"
	"strings"
)


type PasswordItem struct{
	Id string
	Password string
	IsDir bool
}

var PasswordsID map[string]PasswordItem

type Page struct {
	Is_root bool
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

		prefix := os.Getenv("PREFIX")
		uri_params := e.URL.Query()
		
		
		subpath := uri_params["path"]
		
		password_path := prefix
		is_root := true
		if len(subpath) > 0{
			password_path = filepath.Join(prefix, subpath[0])	
			is_root = false
		}
		log.Println(password_path)	
		cmd := exec.Command("ls", password_path)
		output, err := cmd.Output()

		if err != nil {
			log.Printf("cmd.Output() failed with %s\n", err)
		}


		lines := strings.Split(string(output), "\n")
		lines = lines[:len(lines) - 1]
		p := Page{}
		p.Is_root = is_root
		for i := 0; i < len(lines); i++{
			file_p := filepath.Join(password_path, lines[i])
			fileInf, err := os.Stat(file_p)

			if err != nil{
				log.Println(fmt.Sprintf("Couldn't find a file %s", file_p) )
				continue
			}

			passwordID := auth.GenerateChallenge(20)
			p.Passwords = append(p.Passwords, PasswordItem{passwordID, lines[i], fileInf.IsDir()})
			PasswordsID[passwordID] = PasswordItem{passwordID, lines[i], fileInf.IsDir()}
		}
		t.Execute(w, p)
	}
}
