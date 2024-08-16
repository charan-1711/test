package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		userTokenArr := strings.Split(tokenString, " ")
		if len(userTokenArr) < 2 {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(userTokenArr[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		userId := int(claims["userId"].(float64))

		ctx := context.WithValue(r.Context(), "userId", userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
