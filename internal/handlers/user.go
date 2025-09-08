package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Rishit1234567889/baseToAdvGo/internal/dtos"
	"github.com/Rishit1234567889/baseToAdvGo/internal/store"
	"github.com/Rishit1234567889/baseToAdvGo/internal/utils"
	"github.com/Rishit1234567889/baseToAdvGo/internal/validation"
)

// login a user 4.1
func (h *Handler) LoginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//user request aka dto
		var req dtos.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid Request payload")
			return
		}
		if err := validation.Validate(&req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		// fetch user fromdb using the store queries
		user, err := h.Queries.GetUserByEmailOrUsername(ctx, store.GetUserByEmailOrUsernameParams{Username: req.Username})
		if err != nil {
			utils.ResponseWithError(w, http.StatusUnauthorized, "invalid creadential")
			return
		}
		if !utils.ComparePassword(user.Password, req.Password) {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid credential")
			return
		}
		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		token, err := utils.GenerateJWT(int64(user.ID), user.Username, jwtKey)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Error generating a token")
			return
		}
		utils.ResponseWithSuccess(w, http.StatusOK, "Login successful", map[string]string{
			"token": token,
		})
	}
}

func (h *Handler) CreateUserHandler() http.HandlerFunc { //3.*
	return func(w http.ResponseWriter, r *http.Request) {
		// create context
		ctx := r.Context()

		//user request aka dto
		var req dtos.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
			return
		}
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "error while hashing password")
			return
		}
		_, err = h.Queries.CreateUser(ctx, store.CreateUserParams{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
		})

		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "error creating user")
			return
		}
		utils.ResponseWithSuccess(w, http.StatusAccepted, "user created", req.Username)

	}
}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1LCJ1c2VybmFtZSI6InJhYSBzaGEiLCJleHAiOjE3NTc0MjEwMzcsImlzcyI6ImFiYXlvbWkifQ.JXuWde_K0O0QeVcxMUkF7p5gFx018F-4E9gSlA_f2ns
