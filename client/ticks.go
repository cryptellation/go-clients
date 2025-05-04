package client

import (
	"context"

	"github.com/cryptellation/ticks/api"
	"go.temporal.io/sdk/workflow"
)

// ListenToTicks listens to ticks from a specific exchange and trading pair.
func (c client) ListenToTicks(
	ctx context.Context,
	exchange, pair string,
	callback func(ctx workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error,
) error {
	return c.ticks.ListenToTicks(ctx, exchange, pair, callback)
}
