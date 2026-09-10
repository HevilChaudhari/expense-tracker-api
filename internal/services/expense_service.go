package services

import (
	"expense-tracker/internal/models"
	"expense-tracker/internal/repositories"
)

type ExpenseService struct {
	expenseRepository *repositories.ExpenseRepository
}

func NewExpenseService(repository *repositories.ExpenseRepository) *ExpenseService {
	service := ExpenseService{
		expenseRepository: repository,
	}

	return &service
}

func (service *ExpenseService) AddExpense(userID int, newExpense models.Expense) error {
	return service.expenseRepository.Add(userID, newExpense)
}

func (service *ExpenseService) GetAllExpenses(userID int) ([]models.Expense, error) {
	return service.expenseRepository.GetAll(userID)
}

func (service *ExpenseService) GetExpenseByID(userID int, id int) (models.Expense, error) {
	return service.expenseRepository.GetByID(userID, id)
}

func (service *ExpenseService) UpdateExpenseByID(userID int, updatedExpense models.Expense) error {
	return service.expenseRepository.UpdateByID(userID, updatedExpense)
}

func (service *ExpenseService) DeleteExpenseByID(userID int, id int) error {
	return service.expenseRepository.DeleteByID(userID, id)
}

func (service *ExpenseService) GetExpenseSummary(userID int) (models.ExpenseSummary, error) {
	expenses, err := service.expenseRepository.GetAll(userID)

	if err != nil {
		return models.ExpenseSummary{}, err
	}

	summary := models.ExpenseSummary{}

	totalExpenses := len(expenses)

	if totalExpenses == 0 {
		return summary, nil
	}

	var totalAmount float32
	highestAmount := expenses[0].Amount
	lowestAmount := expenses[0].Amount

	for _, expense := range expenses {
		totalAmount += expense.Amount

		if expense.Amount > highestAmount {
			highestAmount = expense.Amount
		}

		if expense.Amount < lowestAmount {
			lowestAmount = expense.Amount
		}
	}

	averageAmount := totalAmount / float32(totalExpenses)

	summary.TotalExpenses = totalExpenses
	summary.TotalAmount = totalAmount
	summary.AverageAmount = averageAmount
	summary.HighestAmount = highestAmount
	summary.LowestAmount = lowestAmount

	return summary, nil
}

