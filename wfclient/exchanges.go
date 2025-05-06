package wfclient

import (
	"github.com/cryptellation/exchanges/api"
	"go.temporal.io/sdk/workflow"
)

// GetExchange gets exchange info from Cryptellation service.
func (c wfClient) GetExchange(
	ctx workflow.Context,
	params api.GetExchangeWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.GetExchangeWorkflowResults, err error) {
	// Set default child workflow options if not provided
	if childWorkflowOptions == nil {
		childWorkflowOptions = &workflow.ChildWorkflowOptions{}
	}

	// Set task queue if not already set
	if childWorkflowOptions.TaskQueue == "" {
		childWorkflowOptions.TaskQueue = api.WorkerTaskQueueName
	}

	// Execute the child workflow with the provided options
	return c.exchanges.GetExchange(ctx, params, childWorkflowOptions)
}
