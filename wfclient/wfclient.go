package wfclient

import (
	candlesticksapi "github.com/cryptellation/candlesticks/api"
	candlesticksclient "github.com/cryptellation/candlesticks/pkg/clients"
	exchangesapi "github.com/cryptellation/exchanges/api"
	exchangesclient "github.com/cryptellation/exchanges/pkg/clients"
	"go.temporal.io/sdk/workflow"
)

// WfClient is a client for the cryptellation exchanges service from a workflow perspective.
type WfClient interface {
	// ListCandlesticks lists candlesticks from Cryptellation service.
	ListCandlesticks(
		ctx workflow.Context,
		params candlesticksapi.ListCandlesticksWorkflowParams,
		childWorkflowOptions *workflow.ChildWorkflowOptions,
	) (result candlesticksapi.ListCandlesticksWorkflowResults, err error)
	// GetExchange calls the exchange get workflow.
	GetExchange(
		ctx workflow.Context,
		params exchangesapi.GetExchangeWorkflowParams,
		childWorkflowOptions *workflow.ChildWorkflowOptions,
	) (result exchangesapi.GetExchangeWorkflowResults, err error)
}

type wfClient struct {
	exchanges    exchangesclient.WfClient
	candlesticks candlesticksclient.WfClient
}

// NewWfClient creates a new workflow client.
// This client is used to call workflows from within other workflows.
// It is not used to call workflows from outside the workflow environment.
func NewWfClient() WfClient {
	return wfClient{
		exchanges: exchangesclient.NewWfClient(),
	}
}
