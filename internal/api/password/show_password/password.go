package show_password

import (
	"encoding/base64"
	"log"
	"net/http"
	"path/filepath"

	// "os"
	"os/exec"
	// "path/filepath"

	show "pass_web/internal/api/show"

	"github.com/gorilla/mux"

	templ "pass_web/internal/api/template"
)

func Handler(w http.ResponseWriter, r *http.Request) {
		t := templ.NewTemplate()
		id := mux.Vars(r)["id"]
		
		passwordItem := show.PasswordsID[id]
		passwordFile := passwordItem.Password
		// prefix := os.Getenv("PREFIX")
		// passwordPath := filepath.Join(prefix, passwordFile)
		//
		passwordPath := passwordItem.Path
		passwordPath = filepath.Join(passwordPath, passwordItem.Password)

		log.Println("Path password ", passwordPath)
		cmd := exec.Command("cat", passwordPath)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Failed to show password: %v", err)
			http.Error(w, "Failed to show password", http.StatusInternalServerError)
			return
		}

		encodedContent := base64.StdEncoding.EncodeToString(output)

		w.Header().Set("Content-Type", "text/html")

		t.Render(w, "password", struct {
			PasswordFile           string
			EncodedContent string
		}{
			PasswordFile :           passwordFile,
			EncodedContent: encodedContent,
		})
	}
