package repositories

import (
	"context"
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

func (repository *ExpenseRepository) Add(newExpense models.Expense) error {

	query := `
	INSERT INTO expenses (title,amount,category) 
	VALUES ($1,$2,$3)
	`

	_, err := repository.pool.Exec(
		context.Background(),
		query,
		newExpense.Title,
		newExpense.Amount,
		newExpense.Category,
	)

	return err
}

func (repository *ExpenseRepository) GetAll() ([]models.Expense, error) {

	rows, err := repository.pool.Query(
		context.Background(),
		`SELECT id, title, amount, category FROM expenses`,
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

func (repository *ExpenseRepository) GetByID(id int) (models.Expense, error) {

	query := `SELECT id,title,amount,category FROM expenses WHERE id = $1`

	row := repository.pool.QueryRow(context.Background(), query, id)

	var expense models.Expense

	err := row.Scan(
		&expense.ID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
	)

	if err != nil {
		return models.Expense{}, err
	}

	return expense, nil
}

func (repository *ExpenseRepository) UpdateByID(newExpense models.Expense) error {

	query := `UPDATE expenses
				SET title = $1,
				amount = $2,
				category = $3
				WHERE id = $4
				`

	_, err := repository.pool.Exec(context.Background(),
		query,
		newExpense.Title,
		newExpense.Amount,
		newExpense.Category,
		newExpense.ID)

	return err
}

func (repository *ExpenseRepository) DeleteByID(id int) error {

	query := `DELETE FROM expenses WHERE id = $1`

	_, err := repository.pool.Exec(context.Background(), query, id)

	return err
}
