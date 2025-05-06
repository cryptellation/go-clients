package wfclient

import (
	"github.com/cryptellation/candlesticks/api"
	"go.temporal.io/sdk/workflow"
)

// ListCandlesticks lists candlesticks from Cryptellation service.
func (c wfClient) ListCandlesticks(
	ctx workflow.Context,
	params api.ListCandlesticksWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.ListCandlesticksWorkflowResults, err error) {
	// Set default child workflow options if not provided
	if childWorkflowOptions == nil {
		childWorkflowOptions = &workflow.ChildWorkflowOptions{}
	}

	// Set task queue if not already set
	if childWorkflowOptions.TaskQueue == "" {
		childWorkflowOptions.TaskQueue = api.WorkerTaskQueueName
	}

	// Execute the child workflow with the provided options
	return c.candlesticks.ListCandlesticks(ctx, params, childWorkflowOptions)
}
