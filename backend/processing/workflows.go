package processing

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type QueryStateResult struct {
  // One of: fetch, save, process, failed, complete
  State    string
  Total    int
  Current  int

}

func ProcessLazadaProducts(ctx workflow.Context, shopID, userID string) (err error) {
  processState := QueryStateResult{
    State: "fetch",
    Total: -1,
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

  for i := 0; i < fetchData.TotalProducts; i++ {
    processState.Current = i
    err = workflow.ExecuteActivity(
      sessionCtx,
      a.SaveLazadaProducts,
      LazadaSaveParam{
        UserID: userID,
        DataKey: fetchData.DataKey,
        ProductIndex: i,
        ProductTotal: fetchData.TotalProducts,
      },
    ).Get(sessionCtx, nil)

    if err != nil {
      processState.State = "failed"
      return err
    }
  }

  processState.State = "process"

  err = workflow.ExecuteActivity(sessionCtx, a.ProcessLazadaProducts, shopID).Get(sessionCtx, nil)

	if err != nil {
    processState.State = "failed"
		return err
	}

  processState.State = "complete"

  return err
}

