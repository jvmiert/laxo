package shop

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/sqlc"
)

var ErrUserNoShops = errors.New("user does have any shops")

type Store interface {
  SaveNewProductToStore(*Product, string) (*sqlc.Product, error)
  GetProductPlatformByProductID(string) (*sqlc.ProductsPlatform, error)
  GetProductPlatformByLazadaID(string) (*sqlc.ProductsPlatform, error)
  CreateProductPlatform(*sqlc.CreateProductPlatformParams) (*sqlc.ProductsPlatform, error)
  UpdateProductToStore(*Product) (*sqlc.Product, error)
  RetrieveShopsByUserID(string) ([]Shop, error)
  GetProductsByShopID(string, int32, int32) ([]sqlc.GetProductsByShopIDRow, error)
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
    return nil, err
  }

  if len(shops) < 1 {
    err := ErrUserNoShops
    return nil, err
  }

  return &shops[0], nil
}

func (s *Service) SaveNewShopToDB(shop *Shop, u string) error {
  _, err := s.store.SaveNewShopToStore(shop, u)

  return err
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

func (s *Service) GetProductListJSON(pp []Product, paginate *Paginate) ([]byte, error) {
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

func (s *Service) GetProductsByUserID(userID string, offset string, limit string) ([]Product, Paginate, error) {
  var pList []Product
  var paginate Paginate

  shops, err := s.store.RetrieveShopsByUserID(userID)
  if err != nil {
    return pList, paginate, fmt.Errorf("RetrieveShopsByUserID: %w", err)
  }

  if len(shops) == 0 {
    return pList, paginate, errors.New("user has not setup any shops yet")
  }

  //@TODO: we don't have an active store logic yet so for now we pick the first
  shopID := shops[0].Model.ID

  offsetI, err := strconv.Atoi(offset)
  if err != nil {
    return pList, paginate, fmt.Errorf("atoi offset: %w", err)
  }

  limitI, err := strconv.Atoi(limit)
  if err != nil {
    return pList, paginate, fmt.Errorf("atoi limit: %w", err)
  }

  pModelList, err := s.store.GetProductsByShopID(shopID, int32(limitI), int32(offsetI))
  if err != nil {
    return pList, paginate, fmt.Errorf("GetProductsByShopID: %w", err)
  }

  total := int64(0)

  if len(pModelList) > 0 {
    total = pModelList[0].Count
  }

  for _, pModel := range pModelList {
    if pModel.ID == "" {
      continue
    }
    mediaListString := string(pModel.MediaIDList)
    mediaList := strings.Split(mediaListString, ",")

    var platformList []ProductPlatformInformation

    lazadaPlatform := ProductPlatformInformation{
      ID: strconv.FormatInt(pModel.LazadaID, 10),
      ProductURL: pModel.LazadaUrl,
      Name: pModel.LazadaName,
      PlatformName: "lazada",
    }

    platformList = append(platformList, lazadaPlatform)

    pList = append(pList, Product{
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

func (s *Service) SaveOrUpdateProductToStore(p *Product, shopID string, lazadaID string) (*Product, error) {
  var platform *sqlc.ProductsPlatform
  var pReturn *Product
  var newModel *sqlc.Product
  var err error

  platform, err = s.GetProductPlatformByLazadaID(lazadaID)
  if err != pgx.ErrNoRows && err != nil {
    return nil, err
  }

  // product was not yet saved
  if err == pgx.ErrNoRows {
    newModel, err = s.store.SaveNewProductToStore(p, shopID)
    if err != nil {
      return nil, err
    }

    param := &sqlc.CreateProductPlatformParams{
      ProductID: newModel.ID,
      ProductsLazadaID: null.StringFrom(lazadaID),
    }
    platform, err = s.store.CreateProductPlatform(param)
    if err != nil {
      return nil, err
    }

    pReturn = &Product{
      Model: newModel,
      PlatformModel: platform,
    }

    return pReturn, nil
  }

  p.Model.ID = platform.ProductID

  newModel, err = s.store.UpdateProductToStore(p)
  if err != nil {
    return nil, err
  }

  pReturn = &Product{
    Model: newModel,
    PlatformModel: platform,
  }

  return pReturn, nil
}
