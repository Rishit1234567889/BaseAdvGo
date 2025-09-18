package routes

import (
	"net/http"

	"github.com/Rishit1234567889/baseToAdvGo/internal/handlers"
	"github.com/Rishit1234567889/baseToAdvGo/internal/middlewares"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {

	mux.HandleFunc("/user/register", handler.CreateUserHandler())
	mux.HandleFunc("/user/login", handler.LoginUserHandler())
	mux.HandleFunc("GET /user/profile", middlewares.AuthMiddleware(http.HandlerFunc(handler.UserProfile())))
	mux.HandleFunc("POST /user/session/logout", middlewares.AuthMiddleware(http.HandlerFunc(handler.LogoutHandler())))

	uploadMux := http.NewServeMux()
	uploadMux.Handle("POST /", middlewares.AuthMiddleware(http.HandlerFunc(handler.UploadProfileImageHandler())))
	mux.Handle("/upload/", http.StripPrefix("/upload", uploadMux))
}
