package workflow

import (
	"errors"
	"time"
)

type BillStatus string

const (
    BillStatusOpen   BillStatus = "open"
    BillStatusClosed BillStatus = "closed"
)

type BillUpdate struct {
	Status BillStatus `json:"status"`
}

type Bill struct {
	CustomerID     string     `json:"customer_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	CurrentCharges []LineItem `json:"current_charges"`
	Currency       string     `json:"currency"`
	Status         BillStatus `json:"status"`
	TotalCharges   float64    `json:"total_charges"`
}

// LineItem represents an individual charge
type LineItem struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
	Currency    string    `json:"currency"`
	Metadata    string    `json:"metadata"`
}

var (
    currencyRates = map[string]float64{
        "USD": 1.0,    // Base currency
        "GEL": 0.37,   // 1 USD = 2.70 GEL (approx.)
        "MYR": 0.22,   // 1 USD = 4.55 MYR (approx.)
        "GBP": 1.25,   // 1 USD = 0.80 GBP (approx.)
    }
)

func(l *LineItem) GetAmountByCurrency(currency string) float64 {
	usdAmount := l.Amount * currencyRates[l.Currency]
	// Convert from USD to the target currency
	return usdAmount / currencyRates[currency]
}

func(b *Bill) Update(update BillUpdate) {
	if update.Status != "" {
		b.Status = update.Status
	}
	if update.Status == BillStatusClosed {
		b.EndDate = time.Now()
	}
}

func (b *Bill) addLineItem(items []LineItem) error {
	if b.Status == BillStatusClosed {
		return errors.New("bill is closed, line items cannot be added")
	}
	existingIDs := make(map[string]bool)
	for _, item := range b.CurrentCharges {
		existingIDs[item.ID] = true
	}

	// Check for duplicates
	for _, newItem := range items {
		if existingIDs[newItem.ID] {
			return errors.New("duplicate line item ID found")
		}
	}

	for _, newItem := range items {
		amount := newItem.GetAmountByCurrency(b.Currency)
		b.TotalCharges += amount
	}
	b.CurrentCharges = append(b.CurrentCharges, items...)
	return nil
}
