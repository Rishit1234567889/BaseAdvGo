package routes

import (
	"github.com/Rishit1234567889/baseToAdvGo/internal/handlers"

	"net/http"
)

func SetupHealthRoute(mux *http.ServeMux, handler *handlers.Handler) { // 1.7 (B) setup Routes
	mux.HandleFunc("/health", handler.HealthHandler())
}
