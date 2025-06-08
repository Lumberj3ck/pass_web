package password

import (
	"encoding/base64"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
	show "pass_web/internal/api/show"

	templ "pass_web/internal/api/template"
)

func Handler(t *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		passwordItem := show.PasswordsID[id]
		passwordFile := passwordItem.Password
		passwordPath := filepath.Join("/root/.password-store", passwordFile)

		cmd := exec.Command("cat", passwordPath)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Failed to list passwords: %v", err)
			http.Error(w, "Failed to list passwords", http.StatusInternalServerError)
			return
		}

		encodedContent := base64.StdEncoding.EncodeToString(output)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<div class="bg-white rounded-lg shadow-md p-6">
				<h2 class="text-xl font-semibold mb-4">` + id + `</h2>
				<div class="bg-gray-100 p-4 rounded">
					<pre class="whitespace-pre-wrap break-all">` + encodedContent + `</pre>
				</div>
			</div>
		`))
	}
}