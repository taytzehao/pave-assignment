package workflow

import (
	"fmt"
    "time"
	"go.temporal.io/sdk/workflow"
)

type CreateBillRequest struct {
    CustomerID string `json:"customer_id"`
	Currency   string `json:"currency"`
}

func BillWorkflow(ctx workflow.Context, createBillRequest CreateBillRequest) (*Bill, error) {
    now := workflow.Now(ctx)
    currency := "USD"
    if createBillRequest.Currency != "" {
        currency = createBillRequest.Currency
    }
    bill := &Bill{
        CustomerID: createBillRequest.CustomerID,
        StartDate: now,
        EndDate: now.AddDate(0, 1, 0),
        Status: BillStatusOpen,
        Currency: currency,
        TotalCharges: 0,
    }
    selector := workflow.NewSelector(ctx)
    if err := workflow.SetQueryHandler(ctx, "getBill", func() (*Bill, error) {
        fmt.Println("LALALA1")
        return bill, nil
    }); err != nil {
        return nil, err
    }

    if err := workflow.SetUpdateHandler(ctx, "updateBill", func(ctx workflow.Context, update BillUpdate) (Bill, error) {
        
        fmt.Printf("%+v\n", update)
        bill.Update(update)
        if update.Status == BillStatusClosed {
            workflow.SignalExternalWorkflow(ctx, workflow.GetInfo(ctx).WorkflowExecution.ID, "", "closeBill", nil)
        }
        return *bill, nil
    }); err != nil {
        fmt.Println("SOMETHING Wreong, err:", err)
        return nil,err
    }

    if err := workflow.SetUpdateHandler(ctx, "addLineItem", func(ctx workflow.Context, lineItems []LineItem) (Bill, error) {
        fmt.Println("bill")
        fmt.Printf("%+v\n", lineItems)

        if err := bill.addLineItem(lineItems); err != nil {
            fmt.Println("Wrong")
            return *bill, err
        }
        
        return *bill, nil
    }); err != nil {
        fmt.Println("SOMETHING Wreong, err:", err)
        return nil, err
    }

    oneMonth := workflow.NewTimer(ctx, bill.EndDate.Sub(now))
    channel := workflow.GetSignalChannel(ctx, "closeBill")

    selector.AddFuture(oneMonth, func(f workflow.Future) {
       bill.Update(BillUpdate{Status: BillStatusClosed})
    })
    selector.AddReceive(channel, func(c workflow.ReceiveChannel, more bool) {
        var closeSignal struct{}
        c.Receive(ctx, &closeSignal)
    })

    
    selector.Select(ctx)

    

	return bill, nil
}