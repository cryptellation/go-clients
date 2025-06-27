package wfclient

import (
	backtestsapi "github.com/cryptellation/backtests/api"
	forwardtestsapi "github.com/cryptellation/forwardtests/api"
	"github.com/cryptellation/runtime"
	"go.temporal.io/sdk/workflow"
)

// SubscribeToPrice subscribes to specific price updates.
func (c wfClient) SubscribeToPrice(ctx workflow.Context, params SubscribeToPriceParams) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue: params.Context.ParentTaskQueue,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	switch params.Context.Mode {
	case runtime.ModeBacktest:
		_, err := c.backtests.SubscribeToPrice(ctx, backtestsapi.SubscribeToPriceWorkflowParams{
			BacktestID: params.Context.ID,
			Exchange:   params.Exchange,
			Pair:       params.Pair,
		})
		return err
	case runtime.ModeForwardtest:
		_, err := c.forwardtests.SubscribeToPrice(ctx, forwardtestsapi.SubscribeToPriceWorkflowParams{
			ForwardtestID: params.Context.ID,
			Exchange:      params.Exchange,
			Pair:          params.Pair,
		})
		return err
	case runtime.ModeLive:
		return ErrNotImplemented
	default:
		return runtime.ErrInvalidMode
	}
}
