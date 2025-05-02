package client

import (
	"context"
	"errors"

	exchangesapi "github.com/cryptellation/exchanges/api"
	exchangesclient "github.com/cryptellation/exchanges/pkg/client"
	temporalclient "go.temporal.io/sdk/client"
	temporalLog "go.temporal.io/sdk/log"
	"golang.org/x/sync/errgroup"
)

// Client is a client for the Cryptellation stack.
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
	temporal struct {
		client temporalclient.Client
		addr   string
		logger temporalLog.Logger
	}

	exchanges exchangesclient.Client
}

// Options is a function that modifies the client configuration.
type Options func(*client)

// WithTemporalAddress sets the address of the temporal server.
// This is used when the temporal client is not provided directly.
func WithTemporalAddress(addr string) func(*client) {
	return func(c *client) {
		c.temporal.addr = addr
	}
}

// WithTemporalClient sets the temporal client directly.
func WithTemporalClient(cl temporalclient.Client) func(*client) {
	return func(c *client) {
		c.temporal.client = cl
	}
}

// WithTemporalLogger sets the logger for the temporal client.
func WithTemporalLogger(logger temporalLog.Logger) func(*client) {
	return func(c *client) {
		c.temporal.logger = logger
	}
}

// New creates a new client to communicate with the Cryptellation stack.
func New(opts ...Options) (Client, error) {
	var c client

	// Apply default options
	c.temporal.logger = &DummyLogger{}

	// Apply options
	for _, opt := range opts {
		opt(&c)
	}

	// Check if either temporal client or address is provided
	switch {
	case c.temporal.client == nil && c.temporal.addr == "":
		return nil, errors.New("temporal client or address must be provided")
	case c.temporal.client != nil && c.temporal.addr != "":
		return nil, errors.New("only one of temporal client or address must be provided")
	case c.temporal.client == nil:
		cl, err := temporalclient.Dial(temporalclient.Options{
			Logger:   c.temporal.logger,
			HostPort: c.temporal.addr,
		})
		if err != nil {
			return nil, err
		}
		c.temporal.client = cl
	}

	// Initialize services
	c.exchanges = exchangesclient.New(c.temporal.client)

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

// GetTemporalClient returns the internal temporal client.
func (c *client) GetTemporalClient() temporalclient.Client {
	return c.temporal.client
}

// Close closes the temporal client if it was created in this package.
// If the client was provided externally, it is the caller's responsibility to close it.
func (c *client) Close() {
	// Close the temporal client if it was created in this package
	if c.temporal.client != nil && c.temporal.addr != "" {
		c.temporal.client.Close()
	}
}
