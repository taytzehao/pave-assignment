package core

import (
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestLineItem_GetAmountByCurrency(t *testing.T) {
	tests := []struct {
		name           string
		lineItem       LineItem
		targetCurrency Currency
		expected       float64
	}{
		{
			name: "USD to GEL",
			lineItem: LineItem{
				Amount:   100,
				Currency: USD,
			},
			targetCurrency: GEL,
			expected:       270.27, // 100 USD = 270.27 GEL (approx.)
		},
		{
			name: "GEL to USD",
			lineItem: LineItem{
				Amount:   270,
				Currency: GEL,
			},
			targetCurrency: USD,
			expected:       99.9, // 270 GEL = 99.9 USD (approx.)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.lineItem.GetAmountByCurrency(tt.targetCurrency)
			if !almostEqual(result, tt.expected, 0.01) {
				t.Errorf("GetAmountByCurrency() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func almostEqual(a, b, epsilon float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}

func TestBillUpdate_Close(t *testing.T) {
	// Create a test bill
	bill := &Bill{
		CustomerID: "test123",
		StartDate:  time.Now().Add(-24 * time.Hour),
		EndDate:    time.Time{}, // Zero value
		Status:     BillStatusOpen,
	}
	
	// Apply the update
	bill.Update(BillStatusClosed)

	// Check if the status was updated correctly
	if bill.Status != BillStatusClosed {
		t.Errorf("Expected status to be %s, got %s", BillStatusClosed, bill.Status)
	}

	// Check if the EndDate was set
	if bill.EndDate.IsZero() {
		t.Error("Expected EndDate to be set, but it's still zero")
	}
}


func TestAddLineItem(t *testing.T) {
	t.Run("Closed Bill", func(t *testing.T) {
		bill := &Bill{
			Status: BillStatusClosed,
		}

		err := bill.AddLineItem(LineItem{
			ID: "1", 
			Description: "Item 1", 
			Amount: 11.4,
			Currency: USD,
			Timestamp: time.Now(),
		},)

		assert.Error(t, err)
		assert.Equal(t, "bill is closed, line items cannot be added", err.Error())
	})

	t.Run("Duplicate Line Item", func(t *testing.T) {
		bill := &Bill{
			Status:         BillStatusOpen,
			Currency:       "USD",
			CurrentCharges: []LineItem{{ID: "1", Description: "Existing Item", Amount: 10.0}},
		}

		err := bill.AddLineItem(LineItem{
			ID: "1", Description: "Duplicate Item", Amount: 15.0,
		},)

		assert.Error(t, err)
		assert.Equal(t, "duplicate line item ID found", err.Error())
	})

	t.Run("Successful Addition", func(t *testing.T) {
		bill := &Bill{
			Status:         BillStatusOpen,
			Currency:       USD,
			CurrentCharges: []LineItem{{ID: "1", Description: "Existing Item", Amount: 10.0, Timestamp: time.Now()}},
			TotalCharges:   10.0,
		}

		newItems := LineItem{
			ID: "2", Description: "New Item 1", Amount: 15.0, Currency: "USD", Timestamp: time.Now(),
		}

		err := bill.AddLineItem(newItems)

		assert.NoError(t, err)
		assert.Len(t, bill.CurrentCharges, 2)
		assert.Equal(t, 25.0, bill.TotalCharges)
	})
}