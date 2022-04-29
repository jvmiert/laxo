package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/guregu/null.v4"
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

  _ = lazada.NewService(store, logger)

  if err = laxo.InitDatabase(dbURI); err != nil {
    logger.Error("Failed to init Database", "uri", dbURI, "error", err)
    return
  }

  clientID := os.Getenv("LAZADA_ID")
  clientSecret := os.Getenv("LAZADA_SECRET")

  accessToken, err := laxo.Queries.GetValidAccessTokenByShopID(
    context.Background(),
    "01G1FZCVYH9J47DB2HZENSBC6E",
  )

  if err != nil {
    logger.Error("Failed getting access token")
    return
  }

  client := lazada.NewClient(clientID, clientSecret, accessToken, logger)

  //response, err := client.QueryProduct(null.NewInt(0, false), null.NewString("fasldfjaffad333", true))
  response, err := client.QueryProduct(null.NewInt(1829634167, true), null.NewString("", false))

  //response, err := client.QueryProducts(lazada.QueryProductsParams{
  //  Limit: 50,
  //  Offset: 0,
  //})

  if err != nil {
    logger.Error("QueryProduct failed", "error", err)
    return
  }

  //logger.Info("Query succeeded", "response", fmt.Sprintf("%+v", response))
  logger.Info("Query succeeded", "variations", response.Data.Variations)

  // FOR SAVING THE REPLY
  //file, _ := json.MarshalIndent(response, "", " ")
  //_ = ioutil.WriteFile("test.json", file, 0644)

  // FOR SAVING PRODUCTS
  //for _, product := range response.Data.Products {
  //  if err = lazadaService.SaveOrUpdateProduct(product, "01G1FZCVYH9J47DB2HZENSBC6E"); err != nil {
  //    logger.Error("SaveOrUpdateProduct return error", "error", err)
  //    return
  //  }
  //}
}
