package middlewares

import (
	"encoding/json"
	"net/http"
	"os"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_TOKEN")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "unauthorized",
				"error":   true,
			})
			return
		}

		if token != os.Getenv("API_TOKEN") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "unauthorized",
				"error":   true,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
