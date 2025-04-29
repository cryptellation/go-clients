package client

import (
	"context"

	exchangesapi "github.com/cryptellation/exchanges/api"
)

// GetExchange retrieves an exchange by name.
func (c client) GetExchange(
	ctx context.Context,
	params exchangesapi.GetExchangeWorkflowParams,
) (exchangesapi.GetExchangeWorkflowResults, error) {
	return c.exchanges.GetExchange(ctx, params)
}

// ListExchanges retrieves a list of exchanges.
func (c client) ListExchanges(
	ctx context.Context,
	params exchangesapi.ListExchangesWorkflowParams,
) (exchangesapi.ListExchangesWorkflowResults, error) {
	return c.exchanges.ListExchanges(ctx, params)
}
