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

func ResponseUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	r, _ := json.Marshal(HttpResponse{
		Message: msg,
		Success: false,
	})
	_, _ = w.Write(r)
}

func ResponseBadRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	r, _ := json.Marshal(HttpResponse{
		Message: msg,
		Success: false,
	})
	_, _ = w.Write(r)
}

func ResponseConflict(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	r, _ := json.Marshal(HttpResponse{
		Message: msg,
		Success: false,
	})
	_, _ = w.Write(r)
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

func ResponseCreated(w http.ResponseWriter, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(HttpResponse{
		Message: msg,
		Success: true,
		Data:    data,
	})
}

func ResponseNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func LoadSession(sess *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sess.LoadAndSave(next)
	}
}
