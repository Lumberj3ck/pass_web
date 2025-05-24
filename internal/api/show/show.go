package show

import (
    "net/http"
    templ "pass_web/internal/api/template"
)

type Page struct{
    Password string
}

func Handler(t *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, e *http.Request) {
		t.Render(w, "index", Page{"YOUUU"})
	}
}
