package repositories

import (
	"context"
	"errors"
	"expense-tracker/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseRepository struct {
	pool *pgxpool.Pool
}

func NewExpenseRepository(pool *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{
		pool: pool,
	}
}

// Add creates a new expense associated with the user
func (repository *ExpenseRepository) Add(userID int, newExpense models.Expense) error {

	query := `
	INSERT INTO expenses (title, amount, category, user_id) 
	VALUES ($1, $2, $3, $4)
	`

	_, err := repository.pool.Exec(
		context.Background(),
		query,
		newExpense.Title,
		newExpense.Amount,
		newExpense.Category,
		userID,
	)

	return err
}

// GetAll returns all expenses belonging to the given user
func (repository *ExpenseRepository) GetAll(userID int) ([]models.Expense, error) {

	query := `SELECT id, title, amount, category, user_id FROM expenses WHERE user_id = $1`

	rows, err := repository.pool.Query(
		context.Background(),
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var expenses []models.Expense

	for rows.Next() {
		var expense models.Expense

		err := rows.Scan(
			&expense.ID,
			&expense.Title,
			&expense.Amount,
			&expense.Category,
			&expense.UserID,
		)

		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetByID returns an expense by ID only if it belongs to the user
func (repository *ExpenseRepository) GetByID(userID int, id int) (models.Expense, error) {

	query := `SELECT id, title, amount, category, user_id FROM expenses WHERE id = $1 AND user_id = $2`

	row := repository.pool.QueryRow(context.Background(), query, id, userID)

	var expense models.Expense

	err := row.Scan(
		&expense.ID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
		&expense.UserID,
	)

	if err != nil {
		return models.Expense{}, err
	}

	return expense, nil
}

// UpdateByID updates an expense only if it belongs to the user
func (repository *ExpenseRepository) UpdateByID(userID int, newExpense models.Expense) error {

	query := `UPDATE expenses
				SET title = $1,
				amount = $2,
				category = $3
				WHERE id = $4 AND user_id = $5
				`

	cmdTag, err := repository.pool.Exec(context.Background(),
		query,
		newExpense.Title,
		newExpense.Amount,
		newExpense.Category,
		newExpense.ID,
		userID)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("expense not found")
	}

	return nil
}

// DeleteByID deletes an expense only if it belongs to the user
func (repository *ExpenseRepository) DeleteByID(userID int, id int) error {

	query := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`

	cmdTag, err := repository.pool.Exec(context.Background(), query, id, userID)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("expense not found")
	}

	return nil
}

