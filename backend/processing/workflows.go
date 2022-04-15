package processing

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func ProcessLazadaProducts(ctx workflow.Context, shopID string) (err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
  sessionCtx := workflow.WithActivityOptions(ctx, ao)

  var dataKey string
  var a *Activities

  err = workflow.ExecuteActivity(sessionCtx, a.FetchLazadaProductsFromAPI, shopID).Get(sessionCtx, &dataKey)

	if err != nil {
		return err
	}

  err = workflow.ExecuteActivity(sessionCtx, a.SaveLazadaProducts, dataKey).Get(sessionCtx, nil)

	if err != nil {
		return err
	}

  err = workflow.ExecuteActivity(sessionCtx, a.ProcessLazadaProducts, shopID).Get(sessionCtx, nil)

	if err != nil {
		return err
	}

  return err
}

