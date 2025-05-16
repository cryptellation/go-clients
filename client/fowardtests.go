package client

import (
	"context"

	"github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/forwardtests/pkg/clients"
)

// NewForwardtest creates a new forwardtest.
func (c client) NewForwardtest(
	ctx context.Context,
	params api.CreateForwardtestWorkflowParams,
) (clients.Forwardtest, error) {
	return c.forwardtests.NewForwardtest(ctx, params)
}

// ListForwardtests lists the forwardtests.
func (c client) ListForwardtests(
	ctx context.Context,
	params api.ListForwardtestsWorkflowParams,
) ([]clients.Forwardtest, error) {
	return c.forwardtests.ListForwardtests(ctx, params)
}
