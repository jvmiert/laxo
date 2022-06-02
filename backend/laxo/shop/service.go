package shop

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/models"
	"laxo.vn/laxo/laxo/sqlc"
)

var ErrUserNoShops = errors.New("user does have any shops")

type Store interface {
  SaveNewProductToStore(*models.Product, string) (*sqlc.Product, error)
  GetProductPlatformByProductID(string) (*sqlc.ProductsPlatform, error)
  GetProductPlatformByLazadaID(string) (*sqlc.ProductsPlatform, error)
  CreateProductPlatform(*sqlc.CreateProductPlatformParams) (*sqlc.ProductsPlatform, error)
  UpdateProductToStore(*models.Product) (*sqlc.Product, error)
  RetrieveShopsByUserID(string) ([]Shop, error)
  GetProductsByShopID(string, int32, int32) ([]sqlc.GetProductsByShopIDRow, error)
  GetProductsByNameOrSKU(string, null.String, null.String, int32, int32) ([]sqlc.GetProductsByNameOrSKURow, error)
  GetProductByID(string) (*sqlc.Product, error)
  GetProductDetails(productID string) (*sqlc.GetProductDetailsByIDRow, error)
  RetrieveShopsPlatformsByUserID(string) ([]sqlc.GetShopsPlatformsByUserIDRow, error)
  SaveNewShopToStore(*Shop, string) (*sqlc.Shop, error)
  GetLazadaPlatformByShopID(string) (*sqlc.PlatformLazada, error)
  SaveNewLazadaPlatform(string, *lazada.AuthResponse) (*sqlc.PlatformLazada, error)
  UpdateLazadaPlatform(string, *lazada.AuthResponse) error
  GetShopByID(string) (*sqlc.Shop, error)
  RetrieveShopsPlatformsByShopID(string) ([]sqlc.ShopsPlatform, error)
  CreateShopsPlatforms(string, string) (sqlc.ShopsPlatform, error)
  RetrieveSpecificPlatformByShopID(string, string) (sqlc.ShopsPlatform, error)
}

type Service struct {
  store       Store
  logger      *laxo.Logger
  server      *laxo.Server
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server) Service {
  return Service {
    store: store,
    logger: logger,
    server: server,
  }
}

func (s *Service) GetActiveShopByUserID(userID string) (*Shop, error) {
  shops, err := s.store.RetrieveShopsByUserID(userID)

  if err != nil {
    return nil, fmt.Errorf("RetrieveShopsByUserID: %w", err)
  }

  if len(shops) < 1 {
    return nil, ErrUserNoShops
  }

  return &shops[0], nil
}

func (s *Service) SaveNewShopToDB(shop *Shop, u string) error {
  _, err := s.store.SaveNewShopToStore(shop, u)

  return fmt.Errorf("SaveNewShopToStore: %w", err)
}

func (s *Service) RetrieveShopsPlatformsByUserID(userID string) ([]sqlc.GetShopsPlatformsByUserIDRow, error) {
  return s.store.RetrieveShopsPlatformsByUserID(userID)
}

func (s *Service) GenerateShopPlatformList(ss []sqlc.GetShopsPlatformsByUserIDRow) ([]byte, error) {
  srList := []ShopReturn{}
  srMap := make(map[string]struct{})
  pMap := make(map[string][]PlatformReturn)

  for _, s := range ss {
    if _, present := srMap[s.ID]; !present {
      var sr ShopReturn
      sr.ID       = s.ID
      sr.UserID   = s.UserID
      sr.Name = s.ShopName

      srList = append(srList, sr)

      srMap[s.ID] = struct{}{}
    }

    created := int64(0)

    if s.PlatformCreated.Valid {
      created = s.PlatformCreated.Time.Unix()
    }

    if s.PlatformID.Valid {
      pMap[s.ID] = append(pMap[s.ID], PlatformReturn{
        ID: s.PlatformID.String,
        Name: s.PlatformName.String,
        Created: created,
      })
    }
  }

  srListU := []ShopReturn{}
  for _, s := range srList {
    if _, present := pMap[s.ID]; present {
      s.Platforms = pMap[s.ID]
      srListU = append(srListU, s)
    } else {
      s.Platforms = make([]PlatformReturn, 0)
      srListU = append(srListU, s)
    }

  }

	shopData := map[string]interface{}{
		"shops": srListU,
		"total": len(srListU),
	}

  bytes, err := json.Marshal(shopData)

  if err != nil {
    return bytes, err
  }
  return bytes, nil
}

func (s *Service) GenerateShopList(ss []Shop) ([]byte, error) {
  sList := []json.RawMessage{}
  for _, s := range ss {
    b, err := s.JSON()
    if err != nil {
      return nil, err
    }
    j := json.RawMessage(b)
    sList = append(sList, j)
  }

	shopData := map[string]interface{}{
		"shops": sList,
		"total": len(sList),
	}

  bytes, err := json.Marshal(shopData)

  if err != nil {
    return bytes, err
  }
  return bytes, nil
}

func (s *Service) GetShopCreateFailure(errs error, printer *message.Printer) validation.Errors {
  errMap := errs.(validation.Errors)

  var errorString string
  for key, err := range errMap {
    ozzoError := err.(validation.Error)
    code := ozzoError.Code()
    params := ozzoError.Params()

    switch code {
    case validation.ErrRequired.Code():
      errorString = printer.Sprintf("cannot be blank")
    case validation.ErrLengthOutOfRange.Code():
      errorString = printer.Sprintf("the length must be between %v and %v", number.Decimal(params["min"]), number.Decimal(params["max"]))
    default:
      errorString = printer.Sprintf("unknown validation error")
    }

    newError := ozzoError.SetMessage(errorString)
    errMap[key] = newError
  }
  return errMap
}

func (s *Service) ValidateNewShop(shop *Shop, printer *message.Printer) error {
  err := shop.ValidateNew()
  if err != nil {
    return s.GetShopCreateFailure(err, printer)
  }

  return nil
}

func (s *Service) GetProductListJSON(pp []models.Product, paginate *Paginate) ([]byte, error) {
  pList := []json.RawMessage{}

  for _, p := range pp {
    b, err := p.JSON()
    if err != nil {
      return nil, err
    }
    j := json.RawMessage(b)
    pList = append(pList, j)
  }

	productData := map[string]interface{}{
		"products": pList,
		"paginate": paginate,
	}

  bytes, err := json.Marshal(productData)
  if err != nil {
    return bytes, err
  }

  return bytes, nil
}

func (s *Service) GetProductByID(productID string) (*models.Product, error) {
  pModel, err := s.store.GetProductDetails(productID)
  if err != nil {
    return nil, fmt.Errorf("GetProductByID: %w", err)
  }

  mediaListString := string(pModel.MediaIDList)
  mediaList := strings.Split(mediaListString, ",")

  var platformList []models.ProductPlatformInformation
  var platformSKU string

  if pModel.LazadaPlatformSku.Valid {
    platformSKU =  strconv.FormatInt(pModel.LazadaPlatformSku.Int64, 10)
  }

  lazadaPlatform := models.ProductPlatformInformation{
    ID: strconv.FormatInt(pModel.LazadaID, 10),
    ProductURL: pModel.LazadaUrl,
    Name: pModel.LazadaName,
    PlatformName: "lazada",
    PlatformSKU: platformSKU,
    SellerSKU: pModel.LazadaSellerSku,
  }

  platformList = append(platformList, lazadaPlatform)

  return &models.Product{
        Model: &sqlc.Product{
          ID: pModel.ID,
          Name: pModel.Name,
          Description: pModel.Description,
          Msku: pModel.Msku,
          SellingPrice: pModel.SellingPrice,
          CostPrice: pModel.CostPrice,
          ShopID: pModel.ShopID,
          MediaID: pModel.MediaID,
          Created: pModel.Created,
          Updated: pModel.Updated,
        },
        MediaList: mediaList,
        Platforms: platformList,
      }, nil
}

func (s *Service) GetProductsByNameOrSKU(userID string, name null.String, msku null.String, offset string, limit string) ([]models.Product, Paginate, error) {
  var pList []models.Product
  var paginate Paginate

  shop, err := s.GetActiveShopByUserID(userID)
  if err != nil {
    return pList, paginate, fmt.Errorf("GetActiveShopByUserID: %w", err)
  }

  offsetI, err := strconv.Atoi(offset)
  if err != nil {
    return pList, paginate, fmt.Errorf("atoi offset: %w", err)
  }

  limitI, err := strconv.Atoi(limit)
  if err != nil {
    return pList, paginate, fmt.Errorf("atoi limit: %w", err)
  }

  nameParsed := null.NewString("", false)
  if name.Valid {
    nameParsed = null.StringFrom("%" + name.String + "%")
  }

  mskuParsed := null.NewString("", false)
  if msku.Valid {
    mskuParsed = null.StringFrom("%" + msku.String + "%")
  }

  pModelList, err := s.store.GetProductsByNameOrSKU(shop.Model.ID, nameParsed, mskuParsed, int32(limitI), int32(offsetI))
  if err != nil {
    return pList, paginate, fmt.Errorf("GetProductsByShopID: %w", err)
  }

  total := int64(0)

  if len(pModelList) > 0 {
    total = pModelList[0].Count
  }

  //@TODO: Handle other platforms
  for _, pModel := range pModelList {
    if pModel.ID == "" {
      continue
    }
    mediaListString := string(pModel.MediaIDList)
    mediaList := strings.Split(mediaListString, ",")

    var platformList []models.ProductPlatformInformation
    var platformSKU string

    if pModel.LazadaPlatformSku.Valid {
      platformSKU =  strconv.FormatInt(pModel.LazadaPlatformSku.Int64, 10)
    }

    lazadaPlatform := models.ProductPlatformInformation{
      ID: strconv.FormatInt(pModel.LazadaID, 10),
      ProductURL: pModel.LazadaUrl,
      Name: pModel.LazadaName,
      PlatformName: "lazada",
      PlatformSKU: platformSKU,
      SellerSKU: pModel.LazadaSellerSku,
    }

    platformList = append(platformList, lazadaPlatform)

    pList = append(pList, models.Product{
      Model: &sqlc.Product{
        ID: pModel.ID,
        Name: pModel.Name,
        Description: pModel.Description,
        Msku: pModel.Msku,
        SellingPrice: pModel.SellingPrice,
        CostPrice: pModel.CostPrice,
        ShopID: pModel.ShopID,
        MediaID: pModel.MediaID,
        Created: pModel.Created,
        Updated: pModel.Updated,
      },
      MediaList: mediaList,
      Platforms: platformList,
    })
  }

  paginate.Total = total
  paginate.Pages = (total + int64(limitI) - 1) / int64(limitI)
  paginate.Limit = int64(limitI)
  paginate.Offset = int64(offsetI)

  return pList, paginate, nil
}

func (s *Service) GetProductsByUserID(userID string, offset string, limit string) ([]models.Product, Paginate, error) {
  var pList []models.Product
  var paginate Paginate

  shop, err := s.GetActiveShopByUserID(userID)
  if err != nil {
    return pList, paginate, fmt.Errorf("GetActiveShopByUserID: %w", err)
  }

  offsetI, err := strconv.Atoi(offset)
  if err != nil {
    return pList, paginate, fmt.Errorf("atoi offset: %w", err)
  }

  limitI, err := strconv.Atoi(limit)
  if err != nil {
    return pList, paginate, fmt.Errorf("atoi limit: %w", err)
  }

  pModelList, err := s.store.GetProductsByShopID(shop.Model.ID, int32(limitI), int32(offsetI))
  if err != nil {
    return pList, paginate, fmt.Errorf("GetProductsByShopID: %w", err)
  }

  total := int64(0)

  if len(pModelList) > 0 {
    total = pModelList[0].Count
  }

  //@TODO: Handle other platforms
  for _, pModel := range pModelList {
    if pModel.ID == "" {
      continue
    }
    mediaListString := string(pModel.MediaIDList)
    mediaList := strings.Split(mediaListString, ",")

    var platformList []models.ProductPlatformInformation
    var platformSKU string

    if pModel.LazadaPlatformSku.Valid {
      platformSKU =  strconv.FormatInt(pModel.LazadaPlatformSku.Int64, 10)
    }

    lazadaPlatform := models.ProductPlatformInformation{
      ID: strconv.FormatInt(pModel.LazadaID, 10),
      ProductURL: pModel.LazadaUrl,
      Name: pModel.LazadaName,
      PlatformName: "lazada",
      PlatformSKU: platformSKU,
      SellerSKU: pModel.LazadaSellerSku,
    }

    platformList = append(platformList, lazadaPlatform)

    pList = append(pList, models.Product{
      Model: &sqlc.Product{
        ID: pModel.ID,
        Name: pModel.Name,
        Description: pModel.Description,
        Msku: pModel.Msku,
        SellingPrice: pModel.SellingPrice,
        CostPrice: pModel.CostPrice,
        ShopID: pModel.ShopID,
        MediaID: pModel.MediaID,
        Created: pModel.Created,
        Updated: pModel.Updated,
      },
      MediaList: mediaList,
      Platforms: platformList,
    })
  }

  paginate.Total = total
  paginate.Pages = (total + int64(limitI) - 1) / int64(limitI)
  paginate.Limit = int64(limitI)
  paginate.Offset = int64(offsetI)

  return pList, paginate, nil
}

func (s *Service) GetProductPlatformByProductID(productID string) (*sqlc.ProductsPlatform, error) {
  return s.store.GetProductPlatformByProductID(productID)
}

func (s *Service) GetProductPlatformByLazadaID(productID string) (*sqlc.ProductsPlatform, error) {
  return s.store.GetProductPlatformByLazadaID(productID)
}

func (s *Service) GetSantizedString(d string) (string) {
  p := bluemonday.StrictPolicy()
  santized := p.Sanitize(d)

  return santized
}

func (s *Service) GetLaxoProductFromLazadaData(p *sqlc.ProductsLazada,
  pAttribute *sqlc.ProductsAttributeLazada, pSKU *sqlc.ProductsSkuLazada) (*models.Product, error) {

  numericPrice := pgtype.Numeric{}
  numericPrice.Set(pSKU.Price.String)

  sanitzedDescription := s.GetSantizedString(pAttribute.Description.String)

  product := &models.Product{
    Model: &sqlc.Product{
      ID: "",
      Name: pAttribute.Name,
      Description : null.StringFrom(sanitzedDescription),
      Msku: null.StringFrom(pSKU.SellerSku),
      SellingPrice: numericPrice,
      CostPrice: pgtype.Numeric{Status: pgtype.Null},
      ShopID: p.ShopID,
    },
  }

  return product, nil
}

func (s *Service) SaveOrUpdateProductToStore(p *models.Product, shopID string, lazadaID string) (*models.Product, error) {
  var platform *sqlc.ProductsPlatform
  var pReturn *models.Product
  var newModel *sqlc.Product
  var err error

  platform, err = s.GetProductPlatformByLazadaID(lazadaID)
  if err != pgx.ErrNoRows && err != nil {
    return nil, fmt.Errorf("GetProductPlatformByLazadaID: %w", err)
  }

  // product was not yet saved
  if err == pgx.ErrNoRows {
    newModel, err = s.store.SaveNewProductToStore(p, shopID)
    if err != nil {
      return nil, fmt.Errorf("SaveNewProductToStore: %w", err)
    }

    param := &sqlc.CreateProductPlatformParams{
      ProductID: newModel.ID,
      ProductsLazadaID: null.StringFrom(lazadaID),
    }
    platform, err = s.store.CreateProductPlatform(param)
    if err != nil {
      return nil, fmt.Errorf("CreateProductPlatform: %w", err)
    }

    pReturn = &models.Product{
      Model: newModel,
      PlatformModel: platform,
    }

    return pReturn, nil
  }

  p.Model.ID = platform.ProductID

  newModel, err = s.store.UpdateProductToStore(p)
  if err != nil {
    return nil, fmt.Errorf("UpdateProductToStore: %w", err)
  }

  pReturn = &models.Product{
    Model: newModel,
    PlatformModel: platform,
  }

  return pReturn, nil
}
