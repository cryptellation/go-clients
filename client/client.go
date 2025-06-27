package client

import (
	"context"
	"errors"

	backtestsapi "github.com/cryptellation/backtests/api"
	"github.com/cryptellation/backtests/pkg/backtest"
	backtestsclient "github.com/cryptellation/backtests/pkg/clients"
	candlesticksapi "github.com/cryptellation/candlesticks/api"
	candlesticksclient "github.com/cryptellation/candlesticks/pkg/clients"
	exchangesapi "github.com/cryptellation/exchanges/api"
	exchangesclient "github.com/cryptellation/exchanges/pkg/clients"
	forwardtestsapi "github.com/cryptellation/forwardtests/api"
	forwardtestsclient "github.com/cryptellation/forwardtests/pkg/clients"
	"github.com/cryptellation/runtime"
	smaapi "github.com/cryptellation/sma/api"
	smaclient "github.com/cryptellation/sma/pkg/clients"
	ticksclient "github.com/cryptellation/ticks/pkg/clients"
	"github.com/google/uuid"
	temporalclient "go.temporal.io/sdk/client"
	temporalLog "go.temporal.io/sdk/log"
	"golang.org/x/sync/errgroup"
)

// Client is a client for the Cryptellation stack.
type Client interface {
	// NewBacktest creates a new backtest.
	NewBacktest(
		ctx context.Context,
		params backtest.Parameters,
		callbacks runtime.Callbacks,
	) (backtestsclient.Backtest, error)
	// GetBacktest gets a backtest.
	GetBacktest(
		ctx context.Context,
		params backtestsapi.GetBacktestWorkflowParams,
	) (backtestsclient.Backtest, error)
	// ListBacktests lists backtests.
	ListBacktests(
		ctx context.Context,
		params backtestsapi.ListBacktestsWorkflowParams,
	) ([]backtestsclient.Backtest, error)

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

	// ListSMA retrieves a list of simple moving averages (SMA) for a specific exchange and trading pair.
	ListSMA(
		ctx context.Context,
		params smaapi.ListWorkflowParams,
	) (res smaapi.ListWorkflowResults, err error)

	// NewForwardtest creates a new forwardtest.
	NewForwardtest(
		ctx context.Context,
		params forwardtestsapi.CreateForwardtestWorkflowParams,
	) (forwardtestsclient.Forwardtest, error)
	// ListForwardtests lists the forwardtests.
	ListForwardtests(
		ctx context.Context,
		params forwardtestsapi.ListForwardtestsWorkflowParams,
	) ([]forwardtestsclient.Forwardtest, error)

	// ListenToTicks listens to ticks from a specific exchange and trading pair.
	ListenToTicks(
		ctx context.Context,
		listener ticksclient.ListenerParams,
		exchange, pair string,
	) error

	// StopListeningToTicks unregisters a callback workflow from ticks for a given exchange and pair.
	StopListeningToTicks(
		ctx context.Context,
		listener uuid.UUID,
		exchange string,
		pair string,
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

	backtests    backtestsclient.Client
	candlesticks candlesticksclient.Client
	exchanges    exchangesclient.Client
	forwardtests forwardtestsclient.Client
	sma          smaclient.Client
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
	c.backtests = backtestsclient.New(c.temporal.client)
	c.candlesticks = candlesticksclient.New(c.temporal.client)
	c.exchanges = exchangesclient.New(c.temporal.client)
	c.forwardtests = forwardtestsclient.New(c.temporal.client)
	c.sma = smaclient.New(c.temporal.client)
	c.ticks = ticksclient.New(c.temporal.client)

	return &c, nil
}

// ServicesInfo retrieves information about the services.
func (c *client) ServicesInfo(ctx context.Context) (map[string]any, error) {
	eg, egCtx := errgroup.WithContext(ctx)
	res := make(map[string]any)
	callbacks := map[string]func(ctx context.Context) (any, error){
		"backtests":    func(ctx context.Context) (any, error) { return c.backtests.Info(ctx) },
		"candlesticks": func(ctx context.Context) (any, error) { return c.candlesticks.Info(ctx) },
		"exchanges":    func(ctx context.Context) (any, error) { return c.exchanges.Info(ctx) },
		"forwardtests": func(ctx context.Context) (any, error) { return c.forwardtests.Info(ctx) },
		"sma":          func(ctx context.Context) (any, error) { return c.sma.Info(ctx) },
		"ticks":        func(ctx context.Context) (any, error) { return c.ticks.Info(ctx) },
	}

	for name, callback := range callbacks {
		eg.Go(func() error {
			r, err := callback(egCtx)
			if err != nil {
				return err
			}

			res[name] = r
			return nil
		})
	}

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
