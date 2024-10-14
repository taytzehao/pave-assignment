package assignment

import (
	"context"
	"fmt"

	"encore.app/assignment/workflow"
	"encore.dev"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	envName           = encore.Meta().Environment.Name
	greetingTaskQueue = envName + "-billing"
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
	w.RegisterActivity(workflow.CreateBill)
	w.RegisterActivity(workflow.UpdateBill)
	w.RegisterActivity(workflow.CreateLineItem)

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
