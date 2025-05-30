package wfclient

import (
	backtestsapi "github.com/cryptellation/backtests/api"
	"github.com/cryptellation/runtime"
	"go.temporal.io/sdk/workflow"
)

// SubscribeToPrice subscribes to specific price updates.
func (c wfClient) SubscribeToPrice(ctx workflow.Context, params SubscribeToPriceParams) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue: params.Run.ParentTaskQueue,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	switch params.Run.Mode {
	case runtime.ModeBacktest:
		_, err := c.backtests.SubscribeToPrice(ctx, backtestsapi.SubscribeToPriceWorkflowParams{
			BacktestID: params.Run.ID,
			Exchange:   params.Exchange,
			Pair:       params.Pair,
		})
		return err
	case runtime.ModeForwardtest:
		return ErrNotImplemented
	case runtime.ModeLive:
		return ErrNotImplemented
	default:
		return runtime.ErrInvalidMode
	}
}
