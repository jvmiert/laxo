package client

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/processing"
)

type Client struct {
  Temporal client.Client
  logger   hclog.Logger
}

func NewClient() (*Client, error) {
  tClient, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
  })

  if err !=nil {
    return nil, err
  }

  c := &Client{
    Temporal: tClient,
  }

  return c, err
}

func (c *Client) StartLazadaPlatformSync(shopID string, userID string) (string, error) {
  workflowOptions := client.StartWorkflowOptions{
    ID:        "product_" + laxo.GetUILD(), //@FIX: This should be moved to a service
    TaskQueue: "product",
  }

	we, err := c.Temporal.ExecuteWorkflow(context.Background(), workflowOptions, processing.ProcessLazadaProducts, shopID, userID)
	if err != nil {
    c.logger.Error("Unable to execute workflow", "error", err)
    return "", err
	}

	c.logger.Info("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
  return we.GetID(), nil
}

