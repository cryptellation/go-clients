package wfclient

import (
	"github.com/cryptellation/exchanges/api"
	exchangesclient "github.com/cryptellation/exchanges/pkg/clients"
	"go.temporal.io/sdk/workflow"
)

// Wf Client is a client for the cryptellation exchanges service from a workflow perspective.
type WfClient interface {
	// GetExchange calls the exchange get workflow.
	GetExchange(
		ctx workflow.Context,
		params api.GetExchangeWorkflowParams,
		childWorkflowOptions *workflow.ChildWorkflowOptions,
	) (result api.GetExchangeWorkflowResults, err error)
}

type wfClient struct {
	exchanges exchangesclient.WfClient
}

// NewWfClient creates a new workflow client.
// This client is used to call workflows from within other workflows.
// It is not used to call workflows from outside the workflow environment.
func NewWfClient() WfClient {
	return wfClient{
		exchanges: exchangesclient.NewWfClient(),
	}
}
