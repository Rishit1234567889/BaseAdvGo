package routes

import (
	"net/http"

	"github.com/Rishit1234567889/baseToAdvGo/internal/handlers"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {

	mux.HandleFunc("/user/register", handler.CreateUserHandler())
}
