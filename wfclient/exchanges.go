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
	return c.exchanges.GetExchange(ctx, params, childWorkflowOptions)
}
