package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
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

  jsonFile, err := os.Open(".\\cmd\\laz_process\\test.json")
  if err != nil {
    logger.Error("Failed to read file", "error", err)
    return
  }

  defer jsonFile.Close()

  byteValue, err := ioutil.ReadAll(jsonFile)
  if err != nil {
    logger.Error("Failed to read bytes from file", "error", err)
    return
  }

  var response lazada.ProductsResponse

  json.Unmarshal(byteValue, &response)

  for i, product := range response.Data.Products {
    iStr := strconv.Itoa(i)
    for j, sku := range response.Data.Products[i].Skus {

      jStr := strconv.Itoa(j)
      logger.Debug("Reading product", "ItemID", product.ItemID, "ShopSku", sku.ShopSku)

      value := gjson.GetBytes(response.RawData, "data.products." + iStr + ".skus." + jStr)
      logger.Debug("Raw sku data", "data", value.String())
    }

    if err = lazadaService.SaveOrUpdateProduct(product, "01G1FZCVYH9J47DB2HZENSBC6E"); err != nil {
      logger.Error("SaveOrUpdateProduct return error", "error", err)
      return
    }
  }
}
