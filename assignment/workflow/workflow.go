package workflow

import (
	"fmt"
	"time"

	"encore.app/assignment/core"
	"encore.app/assignment/interfaces"
	"go.temporal.io/sdk/workflow"
)

func BillWorkflow(ctx workflow.Context, createBillRequest interfaces.CreateBillRequest) (*core.Bill, error) {
	now := workflow.Now(ctx)
	currency := core.USD
	if createBillRequest.Currency != "" {
		currency = createBillRequest.Currency
	}
	bill := &core.Bill{
		ID:           workflow.GetInfo(ctx).WorkflowExecution.ID,
		CustomerID:   createBillRequest.CustomerID,
		StartDate:    now,
		EndDate:      now.AddDate(0, 1, 0),
		Status:       core.BillStatusOpen,
		Currency:     currency,
		TotalCharges: 0,
	}

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5, // Adjust this value as needed
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOptions)

	selector := workflow.NewSelector(ctx)

	if err := workflow.SetQueryHandler(ctx, "getBill", func() (*core.Bill, error) {
		return bill, nil
	}); err != nil {
		return nil, fmt.Errorf("failed to set getBill query handler: %w", err)
	}

	if err := workflow.SetUpdateHandler(ctx, "updateBill", func(ctx workflow.Context, status core.BillStatus) (core.Bill, error) {
		bill.Update(status)
		// workflow.ExecuteActivity(activityCtx, UpdateBill, bill).Get(ctx, nil)
		return *bill, nil
	}); err != nil {
		return nil, err
	}

	if err := workflow.SetUpdateHandler(ctx, "addLineItem", func(ctx workflow.Context, lineItems core.LineItem) (core.Bill, error) {
		if err := bill.AddLineItem(lineItems); err != nil {
			return *bill, err
		}
		workflow.ExecuteActivity(activityCtx, CreateLineItem, &lineItems, bill.ID).Get(ctx, nil)
		return *bill, nil
	}); err != nil {
		return nil, err
	}

	if err := workflow.ExecuteActivity(activityCtx, CreateBill, bill).Get(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to create bill: %w", err)
	}

	oneMonth := workflow.NewTimer(ctx, bill.EndDate.Sub(now))
	channel := workflow.GetSignalChannel(ctx, "closeBill")

	selector.AddFuture(oneMonth, func(f workflow.Future) {
		bill.Update(core.BillStatusClosed)
		workflow.ExecuteActivity(ctx, UpdateBill, bill).Get(ctx, nil)
	})
	selector.AddReceive(channel, func(c workflow.ReceiveChannel, more bool) {
		var closeSignal struct{}
		c.Receive(ctx, &closeSignal)
	})

	selector.Select(ctx)

	return bill, nil
}
