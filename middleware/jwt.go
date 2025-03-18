package middleware

import (
	"fmt"
	"go-crud-database/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.WriteJson(w, http.StatusUnauthorized, "error", nil, "Unauthorized")
			return
		}

		token , err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
			}
			return jwtKey, nil
		})
		
		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
			
		next.ServeHTTP(w, r)
	})

}