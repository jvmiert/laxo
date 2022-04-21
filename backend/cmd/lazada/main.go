package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/sqlc"
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

  //logger.Info("Query succeeded", "response", fmt.Sprintf("%+v", response))
  logger.Info("Query succeeded", "total", response.Data.TotalProducts)

  // @TODO: figure out a way to properly save these products with the null
  //        posibility.
  //
  //        also need to check if the product already exists by checking
  //        the lazada id? maybe make it unique?

  for _, product := range response.Data.Products {
    params := sqlc.CreateLazadaProductParams{
      LazadaID: product.ItemID.Int64,
      LazadaPrimaryCategory: product.PrimaryCategory.Int64,
      Created: product.CreatedTime,
      Updated: product.UpdatedTime,
      Status: product.Status.String,
      SubStatus: product.SubStatus.String,
      ShopID: "01FZQNSR6AQCZ5CZ09QB80AH3C",
    }
    _, err := laxo.Queries.CreateLazadaProduct(context.Background(), params)
    if err != nil {
      logger.Error("Product save error", "error", err)
      return
    }
  }
}
