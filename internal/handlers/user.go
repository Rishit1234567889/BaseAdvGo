package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// "strings"
	"time"

	"github.com/Rishit1234567889/baseToAdvGo/internal/dtos"
	"github.com/Rishit1234567889/baseToAdvGo/internal/middlewares"
	"github.com/Rishit1234567889/baseToAdvGo/internal/store"
	"github.com/Rishit1234567889/baseToAdvGo/internal/utils"
	"github.com/Rishit1234567889/baseToAdvGo/internal/validation"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

// 5.1
// profile.
func (h *Handler) UserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(middlewares.UserClaimsKey).(*utils.Claims)
		if !ok {
			utils.ResponseWithError(w, http.StatusBadRequest, "please login  to continue ")
			return
		}
		userID := claims.UserID

		// check the redis first //7.3
		cacheKey := fmt.Sprintf("user: %d", userID)
		if cached, err := h.Redis.Get(r.Context(), cacheKey).Result(); err == nil {
			var user store.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				utils.ResponseWithSuccess(w, http.StatusOK, "success (from cache/redis)", user)
				return
			}
		}

		//Fallback to DB
		user, err := h.Queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			utils.ResponseWithError(w, http.StatusNotFound, "user not found")
			return
		}

		//set to redis 7.4
		userJSON, _ := json.Marshal(user)
		h.Redis.Set(r.Context(), cacheKey, userJSON, 5*time.Minute)
		utils.ResponseWithSuccess(w, http.StatusOK, "success", user)
	}
}

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

// upload user profile 9.2
func (h *Handler) UploadProfileImageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middlewares.UserClaimsKey).(*utils.Claims)
		if !ok {
			utils.ResponseWithError(w, http.StatusBadRequest, "please login  to continue ")
			return
		}
		userID := claims.UserID
		//upload from the form data
		err := r.ParseMultipartForm(10 << 20) // max file of 10MB
		if err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Error parsing data")
			return
		}
		file, fileHeader, err := r.FormFile("profile_image")
		if err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Error retrieving file")
			return
		}
		defer file.Close()
		cld, err := cloudinary.NewFromParams(
			os.Getenv("CLOUNDINARY_CLOUD_NAME"),
			os.Getenv("CLOUNDINARY_API_KEY"),
			os.Getenv("CLOUDINARY_API_SECRET"),
		)
		if err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Error initiating cloudinary")
			return
		}

		uploadedResult, err := cld.Upload.Upload(r.Context(), file, uploader.UploadParams{
			Folder:   "profile_images",
			PublicID: fileHeader.Filename,
		})

		if err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Error uploading image")
			return
		}

		// commit to db
		h.Queries.CreateUserProfile(r.Context(), store.CreateUserProfileParams{
			UserID:       int32(userID),
			ProfileImage: sql.NullString{String: uploadedResult.SecureURL, Valid: uploadedResult.SecureURL != ""},
		})

		utils.ResponseWithSuccess(w, http.StatusOK, "Image uploaded successfully", uploadedResult.SecureURL)
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

// logout handler 8.0
func (h *Handler) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// extract the jwt claims from the context
		claims, ok := r.Context().Value(middlewares.UserClaimsKey).(*utils.Claims)
		if !ok {
			utils.ResponseWithError(w, http.StatusBadRequest, "Please login to continue")
			return
		}

		// extract the token from the auth header
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		// convert exxireAt to time.Time
		expirationTime := time.Unix(claims.ExpiresAt, 0)
		now := time.Now()
		ttl := expirationTime.Sub(now)
		if ttl <= 0 {
			ttl = 5 * time.Minute // fallback ttl
		}

		err := h.Redis.Set(r.Context(), tokenString, "blacklisted", ttl).Err()
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "failed to blacklist token")
			return
		}

		// clean user session in redis

		userIDstr := fmt.Sprintf("%d", claims.UserID)
		if err := h.cleanUserSession(userIDstr); err != nil {
			fmt.Printf("Error cleaning session for %s : %v\n", userIDstr, err)
			return
		}
		utils.ResponseWithSuccess(w, http.StatusOK, "Logged out successfully", true)
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	// parts := strings.SplitN(authHeader, " ", 2)
	// if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
	// 	return ""
	// }
	return authHeader
}

func (h *Handler) cleanUserSession(userID string) error {

	//session :123:*
	pattern := fmt.Sprintf("session: %s:*", userID)
	//Background context for redis
	ctx := context.Background()
	// scan to iterate over all the keys matching the pattern declared
	iter := h.Redis.Scan(ctx, 0, pattern, 0).Iterator()

	// loop through each key from redis
	for iter.Next(ctx) {

		// delete the key from redis
		err := h.Redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			fmt.Printf("failed to delete session")
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}
