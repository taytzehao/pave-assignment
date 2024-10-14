package core

import (
	"errors"
	"time"
)

type BillStatus string

const (
    BillStatusOpen   BillStatus = "open"
    BillStatusClosed BillStatus = "closed"
)

type Bill struct {
	ID             string     `json:"id"`
	CustomerID     string     `json:"customerID" validate:"required"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	CurrentCharges []LineItem `json:"currentCharges"`
	Currency       Currency     `json:"currency"`
	Status         BillStatus `json:"status"`
	TotalCharges   float64    `json:"totalCharges"`
}

// LineItem represents an individual charge
type LineItem struct {
	ID          string    `json:"id" validate:"required,uuid"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Timestamp   time.Time `json:"timestamp" validate:"required"`
	Currency    Currency    `json:"currency" validate:"required"`
	Metadata    string    `json:"metadata"`
}

type Currency string

const (
    USD Currency = "USD"
    GEL Currency = "GEL"
    MYR Currency = "MYR"
    GBP Currency = "GBP"
)

var (
    currencyRates = map[Currency]float64{
        USD: 1.0,    // Base currency
        GEL: 0.37,   // 1 USD = 2.70 GEL (approx.)
        MYR: 0.22,   // 1 USD = 4.55 MYR (approx.)
        GBP: 1.25,   // 1 USD = 0.80 GBP (approx.)
    }
)

func(l *LineItem) GetAmountByCurrency(currency Currency) float64 {
	usdAmount := l.Amount * currencyRates[l.Currency]
	// Convert from USD to the target currency
	return usdAmount / currencyRates[currency]
}

func(b *Bill) Update(status BillStatus) {
	if status != "" {
		b.Status = status
	}
	if status == BillStatusClosed {
		b.EndDate = time.Now()
	}
}

func (b *Bill) AddLineItem(items LineItem) error {
	if b.Status == BillStatusClosed {
		return errors.New("bill is closed, line items cannot be added")
	}
	existingIDs := make(map[string]bool)
	for _, item := range b.CurrentCharges {
		existingIDs[item.ID] = true
	}

	// Check for duplicates
	
	if existingIDs[items.ID] {
		return errors.New("duplicate line item ID found")
	}
	

	
	amount := items.GetAmountByCurrency(Currency(b.Currency))
	b.TotalCharges += amount
	
	b.CurrentCharges = append(b.CurrentCharges, items)
	return nil
}





