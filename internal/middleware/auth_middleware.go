package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"expense-tracker/internal/auth"
)

// Custom key type for context to avoid collisions
type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		// 1. Check if Authorization header is present
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Check if Bearer prefix is used
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// 3. Extract the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Get JWT secret
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "my_expense_tracker_secret_key"
		}

		// 5. Validate the token
		claims, err := auth.ValidateToken(tokenString, secret)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 6. Save userID into request context so handlers can access it if needed
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// 7. Proceed to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID retrieves the authenticated user's ID from the context
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}


