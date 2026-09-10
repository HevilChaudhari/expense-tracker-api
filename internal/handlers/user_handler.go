package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"expense-tracker/internal/auth"
	"expense-tracker/internal/models"
	"expense-tracker/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(newService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: newService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse returns the generated JWT token alongside user details
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	user, err := handler.userService.RegisterUser(r.Context(), request.Name, request.Email, request.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := handler.userService.LoginUser(
		r.Context(),
		request.Email,
		request.Password,
	)

	if err != nil {
		// Return 404 if the user is not registered
		if err.Error() == "user not found" {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		// Return 401 for invalid password or other auth errors
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my_expense_tracker_secret_key"
	}

	// Generate JWT token for this user
	token, err := auth.GenerateToken(user.ID, secret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send back both token and user details
	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

