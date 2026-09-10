package main

import (
	"fmt"
	"net/http"

	"expense-tracker/internal/database"
	"expense-tracker/internal/handlers"
	"expense-tracker/internal/middleware"
	"expense-tracker/internal/repositories"
	"expense-tracker/internal/services"
)

func main() {

	pool, err := database.Connect()

	if err != nil {
		fmt.Println("Database connection error:", err)
		return
	}

	fmt.Println("Database connected successfully!")
	expenseRepository := repositories.NewExpenseRepository(pool)
	expenseService := services.NewExpenseService(expenseRepository)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	userRepository := repositories.NewUserRepository(pool)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	defer pool.Close()

	// Protected routes (require valid JWT in Authorization header)
	http.Handle("/expenses", middleware.AuthMiddleware(http.HandlerFunc(expenseHandler.ExpenseHandler)))
	http.Handle("/expenses/", middleware.AuthMiddleware(http.HandlerFunc(expenseHandler.ExpenseHandler)))
	http.Handle("/expenses/summary", middleware.AuthMiddleware(http.HandlerFunc(expenseHandler.ExpenseHandler)))

	// Public routes (no authentication required)
	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/login", userHandler.Login)


	fmt.Println("Server running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
