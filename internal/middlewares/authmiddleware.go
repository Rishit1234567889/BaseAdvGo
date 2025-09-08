package middlewares

// 5.0
import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Rishit1234567889/baseToAdvGo/internal/utils"
	"github.com/dgrijalva/jwt-go"
)

// creates a custom type for context key to avoid collision
type contextKey string

// constant used in storing our user claims
const UserClaimsKey contextKey = "claims"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieves the Authorization header from the request (postman)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ResponseWithError(w, http.StatusUnauthorized, "No Token provided")
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		claims := &utils.Claims{}

		// parse the token and also validate it
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// we provide our key from the evniroment variable and validate it again  the token from the request
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil

		})
		//handle the validation error
		if err != nil {
			// handling likely tampered token
			if err == jwt.ErrSignatureInvalid {
				utils.ResponseWithError(w, http.StatusBadRequest, "Invalid token signature")
				return
			}
			// handles any other parsing error e.g expired ,malformed etc
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid token")
			return
		}
		// if token is valid ,store the claims in the request context
		if token.Valid {
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			r = r.WithContext(ctx) // replace reqeust context with the new one
			next.ServeHTTP(w, r)   // calls the next handler ,with the updated request
		} else {
			utils.ResponseWithError(w, http.StatusUnauthorized, "invalid token")
		}
	})
}
