// Service hello implements a simple hello world REST API.
package hello

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"encore.dev"
	"encore.dev/rlog"
	"encore.app/hello/workflow"
)


var (
    envName = encore.Meta().Environment.Name
    greetingTaskQueue = envName + "-greeting"
)

//encore:service
type Service struct {
	client client.Client
	worker worker.Worker
}

func initService() (*Service, error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %v", err)
	}

	w := worker.New(c, "bill", worker.Options{})

	w.RegisterWorkflow(workflow.BillWorkflow)

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %v", err)
	}
	return &Service{client: c, worker: w}, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}



type AddLineItemRequest struct {
	BillID   string              `json:"billID,required"`
    LineItem []workflow.LineItem `json:"line_item"`
}



type GreetResponse struct {
    Greeting string
}

type BillResponse struct {
	Bill workflow.Bill
}
// Welcome to Encore!
// This is a simple "Hello World" project to get you started.
//
// To run it, execute "encore run" in your favorite shell.

// ==================================================================

// This is a public REST API that responds with a personalized greeting.
// Learn more about defining APIs with Encore:
// https://encore.dev/docs/primitives/services-and-apis
//
// To call it, run in your terminal:
//
//	curl http://localhost:4000/hello/World
//


// encore:api public method=POST path=/bill/create
func (s *Service)CreateBill(ctx context.Context, createBillRequest workflow.CreateBillRequest) (*BillResponse, error) {
	options := client.StartWorkflowOptions{
        ID:        "bill-workflow",
        TaskQueue: "bill",
    }
	fmt.Println("createBillRequest", createBillRequest)
    we, err := s.client.ExecuteWorkflow(ctx, options, workflow.BillWorkflow, createBillRequest)
    if err != nil {
        return nil, err
    }
    rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

    // Get the results
    var bill workflow.Bill
    resp, err := s.client.QueryWorkflow(ctx, we.GetID(), "", "getBill")
	if err != nil {
		return nil, err
	}
	err = resp.Get(&bill)
	if err != nil {
		return nil, err
	}
    return &BillResponse{Bill: bill}, nil
}

// encore:api public method=POST path=/lineitem/create
func (s *Service) AddLineItem(ctx context.Context, addLineItemRequest AddLineItemRequest) (*BillResponse, error) {
	fmt.Printf("%+v\n", addLineItemRequest.LineItem)
	handle, err := s.client.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		WorkflowID: addLineItemRequest.BillID,
		UpdateName: "addLineItem",
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args: []interface{}{addLineItemRequest.LineItem},
	})
	if err != nil {
		return nil, err
	}
	var bill workflow.Bill
	err = handle.Get(ctx, &bill)
	if err != nil {
		return nil, err
	}
	return &BillResponse{Bill: bill}, nil
}

// encore:api public method=POST path=/bill/update/:billID
func (s *Service) UpdateBill(ctx context.Context, billID string, updateBillRequest workflow.BillUpdate) (*BillResponse, error) {
	handle, err := s.client.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		WorkflowID: billID,
		UpdateName: "updateBill",
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args: []interface{}{updateBillRequest},
	})
	if err != nil {
		return nil, err
	}
	var bill workflow.Bill
	err = handle.Get(ctx, &bill)
	if err != nil {
		return nil, err
	}
	if updateBillRequest.Status == workflow.BillStatusClosed {
		s.client.SignalWorkflow(ctx, billID, "", "closeBill", nil)
	}
	return &BillResponse{Bill: bill}, nil
}

// encore:api public method=GET path=/bill/:billID
func (s *Service) GetBill(ctx context.Context, billID string) (*BillResponse, error) {
	resp, err := s.client.QueryWorkflow(ctx, billID, "", "getBill")
	if err != nil {
		return nil, err
	}
	var bill workflow.Bill
	err = resp.Get(&bill)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Bill object: %+v\n", bill)
	return &BillResponse{Bill: bill}, nil
}

type Response struct {
	Message string
}

// ==================================================================

// Encore comes with a built-in local development dashboard for
// exploring your API, viewing documentation, debugging with
// distributed tracing, and more:
//
//     http://localhost:9400
//

// ==================================================================

// Next steps
//
// 1. Deploy your application to the cloud
//
//     git add -A .
//     git commit -m 'Commit message'
//     git push encore
//
// 2. To continue exploring Encore, check out one of these topics:
//
// 	  Defining APIs and Services:	 https://encore.dev/docs/primitives/services-and-apis
//    Using SQL databases:  		 https://encore.dev/docs/develop/databases
//    Authenticating users: 		 https://encore.dev/docs/develop/auth
//    Building a Slack bot: 		 https://encore.dev/docs/tutorials/slack-bot
//    Building a REST API:  		 https://encore.dev/docs/tutorials/rest-api
//	  Building an Event-Driven app:  https://encore.dev/docs/tutorials/uptime
