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

func (service *ExpenseService) AddExpense(newExpense models.Expense) error {
	return service.expenseRepository.Add(newExpense)
}

func (service *ExpenseService) GetAllExpenses() ([]models.Expense, error) {
	return service.expenseRepository.GetAll()
}

func (service *ExpenseService) GetExpenseByID(id int) (models.Expense, error) {
	return service.expenseRepository.GetByID(id)
}

func (service *ExpenseService) UpdateExpenseByID(updatedExpense models.Expense) error {
	return service.expenseRepository.UpdateByID(updatedExpense)
}

func (service *ExpenseService) DeleteExpenseByID(id int) error {
	return service.expenseRepository.DeleteByID(id)
}

func (service *ExpenseService) GetExpenseSummary() (models.ExpenseSummary, error) {
	expenses, err := service.expenseRepository.GetAll()

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
