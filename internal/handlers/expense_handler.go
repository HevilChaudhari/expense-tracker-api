package handlers

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"expense-tracker/internal/services"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ExpenseHandlerStruct struct {
	expenseService *services.ExpenseService
}

func NewExpenseHandler(service *services.ExpenseService) *ExpenseHandlerStruct {
	handler := ExpenseHandlerStruct{
		expenseService: service,
	}
	return &handler
}

func (handler *ExpenseHandlerStruct) ExpenseHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handler.getExpense(w, r)
	case http.MethodPost:
		handler.addExpense(w, r)
	case http.MethodPut:
		handler.updateExpense(w, r)
	case http.MethodDelete:
		handler.deleteExpense(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (handler *ExpenseHandlerStruct) getExpense(w http.ResponseWriter, r *http.Request) {

	category := r.URL.Query().Get("category")
	minAmount := r.URL.Query().Get("minAmount")
	maxAmount := r.URL.Query().Get("maxAmount")

	var minAmountVal float32
	var maxAmountVal float32

	if minAmount != "" {
		value, err := strconv.ParseFloat(minAmount, 64)

		if err != nil {
			http.Error(w, "Invalid minAmount", http.StatusBadRequest)
			return
		}

		minAmountVal = float32(value)
	}

	if maxAmount != "" {
		value, err := strconv.ParseFloat(maxAmount, 64)

		if err != nil {
			http.Error(w, "Invalid maxAmount", http.StatusBadRequest)
			return
		}

		maxAmountVal = float32(value)
	}

	// Check if any query filter was provided
	if category != "" || minAmount != "" || maxAmount != "" {

		expenses, err := handler.expenseService.GetAllExpenses()

		if err != nil {
			http.Error(w, "Could not get expenses", http.StatusInternalServerError)
			return
		}

		filteredExpense := []models.Expense{}

		for _, expense := range expenses {

			if category != "" && expense.Category != category {
				continue
			}

			if minAmount != "" && expense.Amount < minAmountVal {
				continue
			}

			if maxAmount != "" && expense.Amount > maxAmountVal {
				continue
			}

			filteredExpense = append(filteredExpense, expense)
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(filteredExpense)

		if err != nil {
			http.Error(w, "Failed to encode expenses", http.StatusInternalServerError)
			return
		}

		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) == 3 {

		if parts[2] == "summary" {

			summary, err := handler.expenseService.GetExpenseSummary()

			if err != nil {
				fmt.Fprintln(w, "Could not get ExpensesSummary")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(summary)

			if err != nil {
				http.Error(w, "Failed to encode summary", http.StatusInternalServerError)
				return
			}

			return
		}

		id, err := strconv.Atoi(parts[2])

		if err != nil {
			http.Error(w, "Invalid expense ID", http.StatusBadRequest)
			return
		}

		var expense models.Expense
		expense, err = handler.expenseService.GetExpenseByID(id)

		if err == nil {
			w.Header().Set("Content-Type", "application/json")

			err = json.NewEncoder(w).Encode(expense)

			if err != nil {
				http.Error(w, "Failed to encode expense", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Expense not found", http.StatusNotFound)
		}

	} else {
		expenses, err := handler.expenseService.GetAllExpenses()

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Could not get expenses", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(expenses)
		if err != nil {
			http.Error(w, "Failed to encode expenses", http.StatusInternalServerError)
			return
		}
	}
}

func (handler *ExpenseHandlerStruct) addExpense(w http.ResponseWriter, r *http.Request) {

	var newExpense models.Expense
	// w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&newExpense)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = handler.expenseService.AddExpense(newExpense)

	if err != nil {
		http.Error(w, "Failed to add expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *ExpenseHandlerStruct) updateExpense(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var updatedExpense models.Expense

	err = json.NewDecoder(r.Body).Decode(&updatedExpense)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedExpense.ID = id

	err = handler.expenseService.UpdateExpenseByID(updatedExpense)

	if err != nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExpense)
}

func (handler *ExpenseHandlerStruct) deleteExpense(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	err = handler.expenseService.DeleteExpenseByID(id)

	if err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Expense not found", http.StatusNotFound)
	}
}
