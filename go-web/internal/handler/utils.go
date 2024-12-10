package handler

import (
	"encoding/json"
	"net/http"

	"github.com/izabelly/go-web/internal/model"
)

func respondJSON(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func handleError(w http.ResponseWriter, status int, message string) {
	body := &model.ResBodyProduct{
		Message: message,
		Data:    nil,
		Error:   true,
	}
	respondJSON(w, status, body)
}
