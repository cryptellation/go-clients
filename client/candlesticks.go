package client

import (
	"context"

	candlesticksapi "github.com/cryptellation/candlesticks/api"
)

// ListCandlesticks calls the candlesticks list workflow.
func (c client) ListCandlesticks(
	ctx context.Context,
	params candlesticksapi.ListCandlesticksWorkflowParams,
) (res candlesticksapi.ListCandlesticksWorkflowResults, err error) {
	return c.candlesticks.ListCandlesticks(ctx, params)
}
