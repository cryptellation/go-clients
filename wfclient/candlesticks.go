package wfclient

import (
	candlesticksapi "github.com/cryptellation/candlesticks/api"
	"go.temporal.io/sdk/workflow"
)

// ListCandlesticks lists candlesticks from Cryptellation service.
func (c wfClient) ListCandlesticks(
	ctx workflow.Context,
	params candlesticksapi.ListCandlesticksWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result candlesticksapi.ListCandlesticksWorkflowResults, err error) {
	return c.candlesticks.ListCandlesticks(ctx, params, childWorkflowOptions)
}
