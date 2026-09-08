package models

type ExpenseSummary struct {
	TotalExpenses int     `json:"totalExpenses"`
	TotalAmount   float32 `json:"totalAmount"`
	AverageAmount float32 `json:"averageAmount"`
	HighestAmount float32 `json:"highestAmount"`
	LowestAmount  float32 `json:"lowestAmount"`
}
