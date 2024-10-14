package workflow

import (
	"context"
	"testing"
	"time"

	"encore.app/assignment/core"
	"encore.app/assignment/interfaces"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
)

func Test_BillWorkflowQuery(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(BillWorkflow)
	createBillRequest := interfaces.CreateBillRequest{CustomerID: "123", Currency: "USD"}
	env.ExecuteWorkflow(BillWorkflow, createBillRequest)

	var result core.Bill
	encodedValue, err := env.QueryWorkflow("getBill")

	require.NoError(t, err, "Failed to query workflow")
	require.NoError(t, encodedValue.Get(&result), "Failed to decode result")

	assert.Equal(t, "123", result.CustomerID, "CustomerID should be 123")
}

func TestBillWorkflow_UpdateBill(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock the CreateBill activity
	env.OnActivity(CreateBill, mock.Anything, mock.Anything).Return(func(ctx context.Context, bill *core.Bill) (*core.Bill, error) {
		// your mock function implementation
		return nil, nil
	})

	// Setup the workflow
	createBillRequest := interfaces.CreateBillRequest{
		CustomerID: "customer123",
		Currency:   "USD",
	}
	billUpdate := core.BillStatusClosed

	updateCallback := &UpdateCallbacks{}
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow("updateBill", "something", updateCallback, billUpdate)
	}, 1*time.Second)

	env.RegisterDelayedCallback(func() {
		var result core.Bill
		encodedValue, err := env.QueryWorkflow("getBill")
		require.NoError(t, err)
		require.NoError(t, encodedValue.Get(&result))

		assert.Equal(t, core.BillStatusClosed, result.Status)
	}, 2*time.Second)
	env.ExecuteWorkflow(BillWorkflow, createBillRequest)

	require.NoError(t, env.GetWorkflowError())

}

func TestBillWorkflow_AddLineItem(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	//env.OnActivity(CreateBill, mock.Anything).Return(nil)
	env.OnActivity(UpdateBill, mock.Anything).Return(nil)

	createBillRequest := interfaces.CreateBillRequest{CustomerID: "456"}

	lineItems := core.LineItem{
		ID:          "1",
		Description: "Item 1",
		Amount:      10.5,
		Currency:    core.USD,
		Timestamp:   time.Now(),
	}

	updateCallback := &UpdateCallbacks{}
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow("addLineItem", "something", updateCallback, &lineItems)

	}, 1*time.Second)

	env.RegisterDelayedCallback(func() {
		var result core.Bill
		encodedValue, err := env.QueryWorkflow("getBill")
		require.NoError(t, err)
		require.NoError(t, encodedValue.Get(&result))

		assert.Equal(t, 1, len(result.CurrentCharges))
		assert.Equal(t, "Item 1", result.CurrentCharges[0].Description)
		assert.Equal(t, 10.5, result.CurrentCharges[0].Amount)

	}, 2*time.Second)

	env.ExecuteWorkflow(BillWorkflow, createBillRequest)

}

func TestBillWorkflow_TimeClose(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.OnActivity(CreateBill, mock.Anything, mock.Anything).Return(func(ctx context.Context, bill *core.Bill) (*core.Bill, error) {
		// your mock function implementation
		return nil, nil
	})
	var result core.Bill
	env.RegisterDelayedCallback(func() {
		encodedValue, _ := env.QueryWorkflow("getBill")
		encodedValue.Get(&result)
		assert.Equal(t, core.BillStatusClosed, result.Status)
	}, 2*30*24*time.Hour)

	createBillRequest := interfaces.CreateBillRequest{CustomerID: "456"}
	env.ExecuteWorkflow(BillWorkflow, createBillRequest)

	require.NoError(t, env.GetWorkflowResult(&result))

	assert.Equal(t, core.BillStatusClosed, result.Status)
}

type UpdateCallbacks struct{}

func (c *UpdateCallbacks) Accept() {
	// Optional: Add any logic for when the update is accepted
}
func (c *UpdateCallbacks) Reject(err error) {
	// Handle rejection (this shouldn't happen in a successful test)

}
func (c *UpdateCallbacks) Complete(success interface{}, err error) {
	// Optional: Add assertions on the success result if needed
}
