package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Rishit1234567889/baseToAdvGo/internal/dtos"
	"github.com/Rishit1234567889/baseToAdvGo/internal/store"
	"github.com/Rishit1234567889/baseToAdvGo/internal/utils"
)

func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create context
		ctx := r.Context()

		//user request aka dto
		var req dtos.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
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
