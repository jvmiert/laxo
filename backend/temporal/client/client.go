package client

import (
	"context"

	"go.temporal.io/sdk/client"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/temporal/lazada"
)

type Client struct {
  temporal client.Client
  logger   *laxo.Logger
}

func NewClient(l *laxo.Logger) (*Client, error) {
  tClient, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
  })

  if err !=nil {
    return nil, err
  }

  c := &Client{
    temporal: tClient,
    logger: l,
  }

  return c, err
}

func (c *Client) Close() {
  c.temporal.Close()
}

func (c *Client) StartLazadaPlatformSync(shopID string, userID string) (string, error) {
  workflowOptions := client.StartWorkflowOptions{
    ID:        "product_" + laxo.GetUILD(),
    TaskQueue: "product",
  }

	we, err := c.temporal.ExecuteWorkflow(context.Background(), workflowOptions, lazada.SyncLazadaPlatform, shopID, userID)
	if err != nil {
    c.logger.Errorw("Unable to execute workflow", "error", err)
    return "", err
	}

	c.logger.Infow("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
  return we.GetID(), nil
}

