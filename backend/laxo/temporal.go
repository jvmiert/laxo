package laxo

import (
	"context"

	"go.temporal.io/sdk/client"
	"laxo.vn/laxo/processing"
)
var TemporalClient client.Client

func InitTemporal() (client.Client, error) {
  c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
  })

  if err !=nil {
    return nil, err
  }

  TemporalClient = c

  return c, err
}

func startTask(shopID, userID string) (string, error) {
  workflowOptions := client.StartWorkflowOptions{
    ID:        "product_" + GetUILD(),
    TaskQueue: "product",
  }

	we, err := TemporalClient.ExecuteWorkflow(context.Background(), workflowOptions, processing.ProcessLazadaProducts, shopID, userID)
	if err != nil {
    Logger.Error("Unable to execute workflow", "error", err)
    return "", err
	}

	Logger.Info("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
  return we.GetID(), nil
}
