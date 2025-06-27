package client

import (
	"context"

	"github.com/cryptellation/ticks/pkg/clients"
	"github.com/google/uuid"
)

// ListenToTicks listens to ticks from a specific exchange and trading pair.
func (c client) ListenToTicks(
	ctx context.Context,
	listener clients.ListenerParams,
	exchange, pair string,
) error {
	return c.ticks.ListenToTicks(ctx, listener, exchange, pair)
}

// StopListeningToTicks unregisters a callback workflow from ticks for a given exchange and pair.
func (c client) StopListeningToTicks(
	ctx context.Context,
	listener uuid.UUID,
	exchange string,
	pair string,
) error {
	return c.ticks.StopListeningToTicks(ctx, listener, exchange, pair)
}
