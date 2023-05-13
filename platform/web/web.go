package web

import (
	"encoding/json"
	"github.com/alexedwards/scs/v2"
	"net/http"
)

type HttpResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
}

func ResponseOK(w http.ResponseWriter, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HttpResponse{
		Message: msg,
		Success: true,
		Data:    data,
	})
}

func RenderErrorPage(w http.ResponseWriter, msg string) {
	TemplateRender(w, "err.page.tmpl", NewTemplateWithErr(msg))
}

func LoadSession(sess *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sess.LoadAndSave(next)
	}
}

func Secured(sess *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := sess.Token(r.Context())
			if token != "" {
				next.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/authors", http.StatusSeeOther)
		})
	}
}

