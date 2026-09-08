package middleware

// import (
// 	"net/http"
// 	"strings"
// 	"expense-tracker/internal/auth"
// )

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		authHeader := r.Header.Get("Authorization")

// 		if authHeader == "" {
// 			http.Error(w, "Authorization header required", http.StatusUnauthorized)
// 			return
// 		}

// 		if !strings.HasPrefix(authHeader, "Bearer ") {
// 			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader,"Bearer ")
// 	})
// }
