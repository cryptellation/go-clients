package client

import (
	"context"

	"github.com/cryptellation/ticks/pkg/clients"
)

// ListenToTicks listens to ticks from a specific exchange and trading pair.
// Note: This method is deprecated. The new ticks client API requires a worker and task queue.
// Use the ticks client directly with a worker for listening to ticks.
func (c client) ListenToTicks(
	ctx context.Context,
	listener clients.ListenerParams,
	exchange, pair string,
) error {
	return c.ticks.ListenToTicks(ctx, listener, exchange, pair)
}
