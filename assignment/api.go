// Service assignment implements a simple assignment world REST API.
package assignment

import (
	"context"
	"fmt"

	"encore.app/assignment/core"
	"encore.app/assignment/interfaces"
	"encore.app/assignment/workflow"
	"encore.dev/rlog"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

// encore:api public method=POST path=/bill
func (s *Service) CreateBill(ctx context.Context, createBillRequest interfaces.CreateBillRequest) (*interfaces.BillResponse, error) {
	options := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "bill",
	}
	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.BillWorkflow, createBillRequest)
	if err != nil {
		return nil, err
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	// Get the results
	var bill core.Bill
	resp, err := s.client.QueryWorkflow(ctx, we.GetID(), "", "getBill")
	if err != nil {
		return nil, err
	}
	err = resp.Get(&bill)
	if err != nil {
		return nil, err
	}
	return &interfaces.BillResponse{Bill: bill}, nil
}

// encore:api public method=POST path=/lineitem
func (s *Service) AddLineItem(ctx context.Context, addLineItemRequest interfaces.AddLineItemRequest) (*interfaces.BillResponse, error) {
	handle, err := s.client.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		WorkflowID:   addLineItemRequest.BillID,
		UpdateName:   "addLineItem",
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args:         []interface{}{addLineItemRequest.LineItem},
	})
	if err != nil {
		return nil, err
	}
	var bill core.Bill
	err = handle.Get(ctx, &bill)
	if err != nil {
		return nil, err
	}
	return &interfaces.BillResponse{Bill: bill}, nil
}

// encore:api public method=PATCH path=/bill/:billID
func (s *Service) UpdateBill(ctx context.Context, billID string, updateBillRequest interfaces.BillUpdateRequest) (*interfaces.BillResponse, error) {
	handle, err := s.client.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		WorkflowID:   billID,
		UpdateName:   "updateBill",
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args:         []interface{}{updateBillRequest.Status},
	})
	if err != nil {
		return nil, err
	}
	var bill core.Bill
	err = handle.Get(ctx, &bill)
	if err != nil {
		return nil, err
	}
	if bill.Status == core.BillStatusClosed {
		s.client.SignalWorkflow(ctx, billID, "", "closeBill", nil)
	}
	return &interfaces.BillResponse{Bill: bill}, nil
}

// encore:api public method=GET path=/bill/:billID
func (s *Service) GetBill(ctx context.Context, billID string) (*interfaces.BillResponse, error) {
	resp, err := s.client.QueryWorkflow(ctx, billID, "", "getBill")
	if err != nil {
		return nil, err
	}
	var bill core.Bill
	err = resp.Get(&bill)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Bill object: %+v\n", bill)
	return &interfaces.BillResponse{Bill: bill}, nil
}
