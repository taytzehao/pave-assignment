package workflow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
)

func setupBillWorkflow(t *testing.T) (*testsuite.WorkflowTestSuite, *testsuite.TestWorkflowEnvironment) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(BillWorkflow)

	return testSuite, env
}

func Test_BillWorkflowQuery(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(BillWorkflow)

	env.ExecuteWorkflow(BillWorkflow, "123")

	var result Bill
	encodedValue, err := env.QueryWorkflow("getBill")

	require.NoError(t, err, "Failed to query workflow")
	require.NoError(t, encodedValue.Get(&result), "Failed to decode result")

	assert.Equal(t, "123", result.CustomerID, "CustomerID should be 123")
	//assert.Equal(t,2,3)
}
type AddLineItemCallbacks struct{}

func (c *AddLineItemCallbacks) Accept() {
	fmt.Println("Accept")
	// Optional: Add any logic for when the update is accepted
}
func (c *AddLineItemCallbacks) Reject(err error) {
	fmt.Println("reject", err)
	// Handle rejection (this shouldn't happen in a successful test)
	
}
func (c *AddLineItemCallbacks) Complete(success interface{}, err error) {
	fmt.Println("complete")
	// Optional: Add assertions on the success result if needed
}
func Test_BillWorkFlow_AddLineItem(t *testing.T) {
	_, env := setupBillWorkflow(t)
	

	env.ExecuteWorkflow(BillWorkflow, "123")
	lineItems := []LineItem{
		{Description: "Item 1", Amount: 11.4},
		{Description: "Item 2", Amount: 10.4},
	}
	//var result Bill
	
	
	
	updateCallbacks := &AddLineItemCallbacks{}
	env.UpdateWorkflow("addLineItem", "uniqueID", updateCallbacks, lineItems)
	assert.Equal(t,2,3)
	// require.NoError(t, err, "Failed to update workflow")

	// assert.Equal(t, 2, len(result.CurrentCharges))
	// assert.Equal(t, "Item 1", result.CurrentCharges[0].Description)
	// assert.Equal(t, 11.4, result.CurrentCharges[0].Amount)
	// assert.Equal(t, "Item 2", result.CurrentCharges[1].Description)
	// assert.Equal(t, 10.4, result.CurrentCharges[1].Amount)
}
