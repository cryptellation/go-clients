package client

import (
	"context"

	"github.com/cryptellation/backtests/api"
	"github.com/cryptellation/backtests/pkg/clients"
)

// NewBacktest creates a new backtest.
func (c client) NewBacktest(
	ctx context.Context,
	params api.CreateBacktestWorkflowParams,
) (clients.Backtest, error) {
	return c.backtests.NewBacktest(ctx, params)
}

// GetBacktest gets a backtest.
func (c client) GetBacktest(
	ctx context.Context,
	params api.GetBacktestWorkflowParams,
) (clients.Backtest, error) {
	return c.backtests.GetBacktest(ctx, params)
}

// ListBacktests lists backtests.
func (c client) ListBacktests(
	ctx context.Context,
	params api.ListBacktestsWorkflowParams,
) ([]clients.Backtest, error) {
	return c.backtests.ListBacktests(ctx, params)
}
