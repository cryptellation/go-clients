package client

import (
	"context"

	"github.com/cryptellation/sma/api"
)

func (c client) ListSMA(
	ctx context.Context,
	params api.ListWorkflowParams,
) (res api.ListWorkflowResults, err error) {
	return c.sma.List(ctx, params)
}
