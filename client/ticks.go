package client

import (
	"context"

	"github.com/cryptellation/ticks/pkg/clients"
)

// ListenToTicks listens to ticks from a specific exchange and trading pair.
func (c client) ListenToTicks(
	ctx context.Context,
	listener clients.ListenerParams,
	exchange, pair string,
) error {
	return c.ticks.ListenToTicks(ctx, listener, exchange, pair)
}
