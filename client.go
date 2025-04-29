package client

import (
	"context"
	"errors"

	exchangesapi "github.com/cryptellation/exchanges/api"
	exchangesclient "github.com/cryptellation/exchanges/pkg/client"
	temporalclient "go.temporal.io/sdk/client"
	"golang.org/x/sync/errgroup"
)

type Client interface {
	// GetExchange retrieves an exchange by name.
	GetExchange(
		ctx context.Context,
		params exchangesapi.GetExchangeWorkflowParams,
	) (exchangesapi.GetExchangeWorkflowResults, error)
	// ListExchanges retrieves a list of exchanges.
	ListExchanges(
		ctx context.Context,
		params exchangesapi.ListExchangesWorkflowParams,
	) (exchangesapi.ListExchangesWorkflowResults, error)

	// ServicesInfo retrieves information about the services.
	ServicesInfo(ctx context.Context) (map[string]any, error)

	GetTemporalClient() temporalclient.Client
	Close()
}

type client struct {
	temporal     temporalclient.Client
	temporalAddr string

	exchanges exchangesclient.Client
}

type Options func(*client)

func WithTemporalAddress(addr string) func(*client) {
	return func(c *client) {
		c.temporalAddr = addr
	}
}

func WithTemporalClient(cl temporalclient.Client) func(*client) {
	return func(c *client) {
		c.temporal = cl
	}
}

func New(opts ...Options) (Client, error) {
	var c client

	// Apply options
	for _, opt := range opts {
		opt(&c)
	}

	// Check if either temporal client or address is provided
	if c.temporal == nil && c.temporalAddr == "" {
		return nil, errors.New("temporal client or address must be provided")
	} else if c.temporal == nil {
		cl, err := temporalclient.Dial(temporalclient.Options{
			HostPort: c.temporalAddr,
		})
		if err != nil {
			return nil, err
		}
		c.temporal = cl
	}

	// Initialize services
	c.exchanges = exchangesclient.New(c.temporal)

	return &c, nil
}

// ServicesInfo retrieves information about the services.
func (c *client) ServicesInfo(ctx context.Context) (map[string]any, error) {
	eg, egCtx := errgroup.WithContext(ctx)
	res := make(map[string]any)

	eg.Go(func() error {
		r, err := c.exchanges.Info(egCtx)
		if err != nil {
			return err
		}

		res["exchanges"] = r
		return nil
	})

	return res, eg.Wait()
}

func (c *client) GetTemporalClient() temporalclient.Client {
	return c.temporal
}

func (c *client) Close() {
	// Close the temporal client if it was created in this package
	if c.temporal != nil && c.temporalAddr != "" {
		c.temporal.Close()
	}
}
