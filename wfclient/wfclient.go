package wfclient

import (
	"errors"

	backtestsclient "github.com/cryptellation/backtests/pkg/clients"
	candlesticksapi "github.com/cryptellation/candlesticks/api"
	candlesticksclient "github.com/cryptellation/candlesticks/pkg/clients"
	exchangesapi "github.com/cryptellation/exchanges/api"
	exchangesclient "github.com/cryptellation/exchanges/pkg/clients"
	forwardtestsclient "github.com/cryptellation/forwardtests/pkg/clients"
	"github.com/cryptellation/runtime"
	"go.temporal.io/sdk/workflow"
)

var (
	// ErrNotImplemented is returned when the function is not implemented.
	ErrNotImplemented = errors.New("not implemented")
)

// WfClient is a client for the cryptellation exchanges service from a workflow perspective.
type WfClient interface {
	// SubscribeToPrice subscribes to specific price updates.
	SubscribeToPrice(
		ctx workflow.Context,
		params SubscribeToPriceParams,
	) error

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

// SubscribeToPriceParams is the parameters to subscribe to price updates.
type SubscribeToPriceParams struct {
	Context  runtime.Context
	Exchange string
	Pair     string
}

type wfClient struct {
	backtests    backtestsclient.WfClient
	exchanges    exchangesclient.WfClient
	candlesticks candlesticksclient.WfClient
	forwardtests forwardtestsclient.WfClient
}

// NewWfClient creates a new workflow client.
// This client is used to call workflows from within other workflows.
// It is not used to call workflows from outside the workflow environment.
func NewWfClient() WfClient {
	return wfClient{
		exchanges:    exchangesclient.NewWfClient(),
		forwardtests: forwardtestsclient.NewWfClient(),
	}
}
