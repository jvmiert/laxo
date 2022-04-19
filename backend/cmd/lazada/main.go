package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
)

func main() {
  logger, _ := laxo.InitConfig(false)

  if err := godotenv.Load(".env"); err != nil {
    logger.Error("Failed to load .env file")
    return
  }

  dbURI := os.Getenv("POSTGRESQL_URL")

  if err := laxo.InitDatabase(dbURI); err != nil {
    logger.Error("Failed to init Database", "uri", dbURI, "error", err)
    return
  }

  clientID := os.Getenv("LAZADA_ID")
  clientSecret := os.Getenv("LAZADA_SECRET")

  accessToken, err := laxo.Queries.GetValidAccessTokenByShopID(
    context.Background(),
    "01FZQNSR6AQCZ5CZ09QB80AH3C",
  )

  if err != nil {
    logger.Error("Failed getting access token")
    return
  }

  client := lazada.NewClient(clientID, clientSecret, accessToken, logger)

  response, err := client.QueryProducts(lazada.QueryProductsParams{
    Limit: 50,
    Offset: 0,
  })

  if err != nil {
    logger.Error("QueryPructs failed", "error", err)
    return
  }

  logger.Info("Query succeeded", "response", fmt.Sprintf("%+v", response))
}
