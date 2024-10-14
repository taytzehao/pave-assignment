package interfaces

import (
	"encore.app/assignment/core"
)

type BillUpdateRequest struct {
	Status core.BillStatus `json:"status"`
}

type CreateBillRequest struct {
	CustomerID string        `json:"customerID" validate:"required"`
	Currency   core.Currency `json:"currency"`
}

type AddLineItemRequest struct {
	BillID   string        `json:"billID" validate:"required,uuid"`
	LineItem core.LineItem `json:"lineItem" validate:"required"`
}

type BillResponse struct {
	Bill core.Bill `json:"bill"`
}
