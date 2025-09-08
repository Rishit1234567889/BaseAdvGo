package routes

import (
	"net/http"

	"github.com/Rishit1234567889/baseToAdvGo/internal/handlers"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) { // 1.7(C)
	SetupHealthRoute(mux, handler)
	SetupUserRoutes(mux, handler)

}
