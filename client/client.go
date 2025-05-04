package client

import (
	"context"
	"errors"

	candlesticksapi "github.com/cryptellation/candlesticks/api"
	candlesticksclient "github.com/cryptellation/candlesticks/pkg/clients"
	exchangesapi "github.com/cryptellation/exchanges/api"
	exchangesclient "github.com/cryptellation/exchanges/pkg/clients"
	"github.com/cryptellation/ticks/api"
	ticksclient "github.com/cryptellation/ticks/pkg/clients"
	temporalclient "go.temporal.io/sdk/client"
	temporalLog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/sync/errgroup"
)

// Client is a client for the Cryptellation stack.
type Client interface {
	// ListCandlesticks calls the candlesticks list workflow.
	ListCandlesticks(
		ctx context.Context,
		params candlesticksapi.ListCandlesticksWorkflowParams,
	) (res candlesticksapi.ListCandlesticksWorkflowResults, err error)
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
	// ListenToTicks listens to ticks from a specific exchange and trading pair.
	ListenToTicks(
		ctx context.Context,
		exchange, pair string,
		callback func(ctx workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error,
	) error

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

	exchanges    exchangesclient.Client
	candlesticks candlesticksclient.Client
	ticks        ticksclient.Client
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
	c.candlesticks = candlesticksclient.New(c.temporal.client)
	c.ticks = ticksclient.New(c.temporal.client)

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

	eg.Go(func() error {
		r, err := c.candlesticks.Info(egCtx)
		if err != nil {
			return err
		}

		res["candlesticks"] = r
		return nil
	})

	eg.Go(func() error {
		r, err := c.ticks.Info(egCtx)
		if err != nil {
			return err
		}

		res["ticks"] = r
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
