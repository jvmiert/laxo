package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/store"
)

func main() {
  logger, _ := laxo.InitConfig(false)

  if err := godotenv.Load(".env"); err != nil {
    logger.Error("Failed to load .env file")
    return
  }

  dbURI := os.Getenv("POSTGRESQL_URL")

  store, err := store.NewStore(dbURI, logger)

  if err != nil {
    logger.Error("Failed to create new store", "error", err)
    return
  }

  lazadaService := lazada.NewService(store, logger)

  if err = laxo.InitDatabase(dbURI); err != nil {
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

  //logger.Info("Query succeeded", "response", fmt.Sprintf("%+v", response))
  logger.Info("Query succeeded", "total", response.Data.TotalProducts)

  for _, product := range response.Data.Products {
    if err = lazadaService.SaveOrUpdateProduct(product, "01FZQNSR6AQCZ5CZ09QB80AH3C"); err != nil {
      logger.Error("SaveOrUpdateProduct return error", "error", err)
      return
    }
  }
}
