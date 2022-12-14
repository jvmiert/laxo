package lazada

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type QueryStateResult struct {
	// One of: fetch, save, process, failed, complete
	State   string
	Total   int64
	Current int64
}

func SyncLazadaPlatform(ctx workflow.Context, shopID string, userID string, overwrite bool) (err error) {
	processState := QueryStateResult{
		State:   "fetch",
		Total:   -1,
		Current: -1,
	}

	queryType := "current_state"

	err = workflow.SetQueryHandler(ctx, queryType, func() (QueryStateResult, error) {
		return processState, nil
	})

	if err != nil {
		processState.State = "failed"
		return err
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Minute,
			MaximumAttempts:    2,
		},
	}
	sessionCtx := workflow.WithActivityOptions(ctx, ao)

	var fetchData LazadaFetchResult
	var a *Activities

	err = workflow.ExecuteActivity(sessionCtx, a.FetchLazadaProductsFromAPI, shopID, userID).Get(sessionCtx, &fetchData)

	if err != nil {
		processState.State = "failed"
		return err
	}

	processState.State = "save"
	processState.Total = fetchData.TotalProducts

	for i := int64(0); i < fetchData.TotalProducts; i++ {
		processState.Current = i
		err = workflow.ExecuteActivity(
			sessionCtx,
			a.SaveLazadaProducts,
			LazadaSaveParam{
				UserID:       userID,
				ShopID:       shopID,
				DataKey:      fetchData.DataKey,
				ProductIndex: i,
				ProductTotal: fetchData.TotalProducts,
				Overwrite:    overwrite,
			},
		).Get(sessionCtx, nil)

		if err != nil {
			processState.State = "failed"
			return err
		}
	}

	err = workflow.ExecuteActivity(sessionCtx, a.CompleteLazadaProducts, userID, fetchData.DataKey).Get(sessionCtx, nil)
	if err != nil {
		processState.State = "failed"
		return err
	}

	processState.State = "complete"

	return err
}
