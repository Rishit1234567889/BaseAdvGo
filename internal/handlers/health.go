package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) HealthHandler() http.HandlerFunc { // 1.7 setup Handlers
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message": "Server is Okay",
		}
		json.NewEncoder(w).Encode(response)
	}
}
