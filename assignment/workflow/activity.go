package workflow

import (
	"context"
	"fmt"
	"encore.app/assignment/database"
	"encore.app/assignment/core"
)

// CreateBill inserts a new bill into the database
func CreateBill(ctx context.Context, bill *core.Bill) (*core.Bill, error) {
	fmt.Println("CREATEBILL")
	query := `
		INSERT INTO Bills (id, customer_id, start_date, end_date, currency, status, total_charges)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, customer_id, start_date, end_date, currency, status, total_charges
	`
	fmt.Println("bill", bill)
	createdBill := &core.Bill{}
	err := database.DB.QueryRow(ctx, query,
		bill.ID,
		bill.CustomerID,
		bill.StartDate.Format("2006-01-02T15:04:05Z07:00"),
		bill.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		bill.Currency,
		bill.Status,
		bill.TotalCharges,
	).Scan(
		&createdBill.ID,
		&createdBill.CustomerID,
		&createdBill.StartDate,
		&createdBill.EndDate,
		&createdBill.Currency,
		&createdBill.Status,
		&createdBill.TotalCharges,
	)
	fmt.Println("createdBill", createdBill)
	if err != nil {
		return nil, err
	}

	return createdBill, nil
}

// UpdateBill updates an existing bill in the database
func UpdateBill(ctx context.Context, bill *core.Bill) (*core.Bill, error) {
	query := `
		UPDATE bills
		SET start_date = $2, end_date = $3, currency = $4, status = $5, total_charges = $6
		WHERE id = $1
		RETURNING id, customer_id, start_date, end_date, currency, status, total_charges
	`
	updatedBill := &core.Bill{}
	err := database.DB.QueryRow(ctx, query,
		bill.ID,
		bill.StartDate,
		bill.EndDate,
		bill.Currency,
		bill.Status,
		bill.TotalCharges,
	).Scan(
		&updatedBill.ID,
		&updatedBill.CustomerID,
		&updatedBill.StartDate,
		&updatedBill.EndDate,
		&updatedBill.Currency,
		&updatedBill.Status,
		&updatedBill.TotalCharges,
	)
	if err != nil {
		return nil, err
	}
	return updatedBill, nil
}

func CreateLineItem(ctx context.Context, lineItem *core.LineItem, billID string) error {
	fmt.Println("CREATELINEITEM############")
	query := `
		INSERT INTO lineitems (id, bill_id, description, amount, timestamp, currency, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := database.DB.Exec(ctx, query,
		lineItem.ID,
		billID,
		lineItem.Description,
		lineItem.Amount,
		lineItem.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		lineItem.Currency,
		lineItem.Metadata,
	)
	return err
}
