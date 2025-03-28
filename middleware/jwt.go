package middleware

import (
	"context"
	"go-crud-database/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Unauthorized")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Invalid auth token")
			return
		}

		tokenString := tokenParts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		userIdFloat, ok := claims["userId"].(float64)
		if !ok {
			utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Invalid auth token")
			return
		}
		userId := int(userIdFloat)

		isAdmin, ok := claims["isAdmin"].(bool)
		if !ok {
			isAdmin = false
		}

		// Save user data in request context
		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "isAdmin", isAdmin)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}
