package shop

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"math/big"
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
var ErrProductNotOwned = errors.New("shop does not own this product")
var ErrProductNotFound = errors.New("no product found")

type Store interface {
	SaveNewProductToStore(p *models.Product, shopID string) (*sqlc.Product, error)
	GetProductPlatformByProductID(string) (*sqlc.ProductsPlatform, error)
	GetProductPlatformByLazadaID(string) (*sqlc.ProductsPlatform, error)
	CreateProductPlatform(*sqlc.CreateProductPlatformParams) (*sqlc.ProductsPlatform, error)
	UpdateProductToStore(*models.Product) (*sqlc.Product, error)
	RetrieveShopsByUserID(string) ([]Shop, error)
	GetProductsByShopID(string, int32, int32) ([]sqlc.GetProductsByShopIDRow, error)
	GetProductsByNameOrSKU(string, null.String, null.String, int32, int32) ([]sqlc.GetProductsByNameOrSKURow, error)
	GetProductByID(string) (*sqlc.Product, error)
	GetProductDetails(productID string, shopID string) (*sqlc.GetProductDetailsByIDRow, error)
	GetProductAssetsByProductID(productID string, shopID string) ([]sqlc.GetProductAssetsByProductIDRow, error)
	RetrieveShopsPlatformsByUserID(string) ([]sqlc.GetShopsPlatformsByUserIDRow, error)
	SaveNewShopToStore(*Shop, string) (*sqlc.Shop, error)
	GetLazadaPlatformByShopID(string) (*sqlc.PlatformLazada, error)
	SaveNewLazadaPlatform(string, *lazada.AuthResponse) (*sqlc.PlatformLazada, error)
	UpdateLazadaPlatform(string, *lazada.AuthResponse) error
	GetShopByID(string) (*sqlc.Shop, error)
	RetrieveShopsPlatformsByShopID(string) ([]sqlc.ShopsPlatform, error)
	CreateShopsPlatforms(string, string) (sqlc.ShopsPlatform, error)
	RetrieveSpecificPlatformByShopID(string, string) (sqlc.ShopsPlatform, error)
	GetAssetByOriginalName(originalName string, shopID string) (*sqlc.Asset, error)
	CheckProductOwner(productID string, shopID string) (string, error)
	UpdateLazadaProductPlatformSync(productID string, state bool) error
	UpdateProductImageOrderRequest(productID string, request *models.ProductImageOrderRequest) error
	GetProductByProductMSKU(Msku string, shopID string) (*sqlc.Product, error)
	CreateProductMedia(assetID string, productID string, status string, order int64) (*sqlc.ProductsMedia, error)
}

type Service struct {
	store  Store
	logger *laxo.Logger
	server *laxo.Server
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server) Service {
	return Service{
		store:  store,
		logger: logger,
		server: server,
	}
}

func (s *Service) CreateNewProduct(p *models.NewProductRequest, shopID string) (*sqlc.Product, error) {
	sellingNumeric := pgtype.Numeric{}
	sellingNumeric.Set(p.SellingPrice)

	costNumeric := pgtype.Numeric{}
	costNumeric.Set(p.CostPrice)

	parsedDescription := s.ParseSlateEmptyChildren(p.Description)

	// @TODO: Should filter out empty text properly

	var stringDescription string

	if len(parsedDescription) != 0 {
		bytes, err := json.Marshal(parsedDescription)
		if err != nil {
			return nil, err
		}
		stringDescription = string(bytes)
	}

	s.logger.Debugw("parsedDescription", "len", len(parsedDescription), "value", parsedDescription)

	pModel := models.Product{
		Model: &sqlc.Product{
			Name: p.Name,
			Msku: null.StringFrom(p.Msku),

			// @TODO: add the non-rich description
			//Description:

			DescriptionSlate: null.StringFrom(stringDescription),
			SellingPrice:     sellingNumeric,
			CostPrice:        costNumeric,
			ShopID:           shopID,
		},
	}

	rModel, err := s.store.SaveNewProductToStore(&pModel, shopID)
	if err != nil {
		return nil, fmt.Errorf("SaveNewProductToStore: %w", err)
	}

	for i, a := range p.Assets {
		_, err := s.store.CreateProductMedia(a.ID, rModel.ID, "active", int64(i))
		if err != nil {
			s.server.Logger.Errorw("CreateProductMedia", "error", err)
			continue
		}
	}

	return rModel, nil
}

// Validates a new product from the product creation frontend form
func (s *Service) ValidateNewProductRequest(p *models.NewProductRequest, shopID string) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required, validation.Length(4, 0)),
		validation.Field(&p.Msku, validation.Required, validation.Length(4, 1024)),
		validation.Field(&p.SellingPrice, validation.Required, validation.Min(1)),
		validation.Field(&p.CostPrice, validation.Min(1)),
	)
	if err != nil {
		return err
	}

	// checking if merchant SKU already exists
	_, err = s.store.GetProductByProductMSKU(p.Msku, shopID)

	// Msku doesn't exist which is what we want
	if err == pgx.ErrNoRows {
		return nil
	}

	// We found a product with the merchant SKU, fail validation
	if err == nil {
		err = validation.Errors{
			"msku": validation.NewError(
				"already_exists",
				"already_exists"),
		}
		return err
	}

	// A non-related error occured
	if err != nil {
		return err
	}

	return err
}

func (s *Service) UpdateProductImageOrderRequest(productID string, request *models.ProductImageOrderRequest) error {
	return s.store.UpdateProductImageOrderRequest(productID, request)
}

func (s *Service) ValidateProductImageOrderRequest(request *models.ProductImageOrderRequest) error {
	return validation.Validate(request.Assets)
}

// Allows the user to change whether the product's data is synced back to the platform
func (s *Service) ChangedProductPlatformSync(request *models.ProductChangedSyncRequest, productID string) error {
	if request.Platform == "lazada" {
		return s.store.UpdateLazadaProductPlatformSync(productID, request.State)
	}

	return errors.New("unknown platform")
}

// This function will check whether the store owns this product. Used before making changes to a product
// Returns an error if not owned, nil when owned.
func (s *Service) ProductIsOwnedByStore(productID string, shopID string) error {
	_, err := s.store.CheckProductOwner(productID, shopID)
	if err != pgx.ErrNoRows && err != nil {
		return fmt.Errorf("CheckProductOwner: %w", err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrProductNotOwned
	}

	return nil
}

func (s *Service) UpdateProductFromRequest(r *models.ProductDetailPostRequest, printer *message.Printer, shopID string, productID string) error {
	sellingNumeric := pgtype.Numeric{}
	sellingNumeric.Set(r.SellingPrice)

	costNumeric := pgtype.Numeric{}
	costNumeric.Set(r.CostPrice)

	parsedDescription := s.ParseSlateEmptyChildren(r.Description)
	bytes, err := json.Marshal(parsedDescription)
	if err != nil {
		err = validation.Errors{
			"description": validation.NewError(
				"invalid_description",
				printer.Sprintf("description is not valid")),
		}
		return err
	}

	m := &models.Product{
		Model: &sqlc.Product{
			Name:             r.Name,
			DescriptionSlate: null.StringFrom(string(bytes)),
			SellingPrice:     sellingNumeric,
			CostPrice:        costNumeric,
			ShopID:           shopID,
			ID:               productID,
		},
	}

	_, err = s.store.UpdateProductToStore(m)
	if errors.Is(err, pgx.ErrNoRows) {
		errReturn := validation.Errors{
			"generalError": validation.NewError(
				"product_not_found",
				printer.Sprintf("could not find your product")),
		}

		return errReturn
	}
	if err != nil {
		s.server.Logger.Errorw("UpdateProductFromRequest", "error", err)
		return validation.Errors{
			"generalError": validation.NewError(
				"general_error",
				printer.Sprintf("something went wrong, please try again later")),
		}
	}

	return err
}

func (s *Service) ValidateProductDetails(p *models.ProductDetailPostRequest, printer *message.Printer) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required, validation.Length(4, 0)),
		validation.Field(&p.Msku, validation.Required, validation.Length(4, 1024)),
		validation.Field(&p.SellingPrice, validation.Required, validation.Min(1)),
		validation.Field(&p.CostPrice, validation.Min(1)),
	)

	if err != nil {
		return GetProductDetailFailure(err, printer)
	}

	return nil
}

func (s *Service) GetShopByID(shopID string) (*Shop, error) {
	shop, err := s.store.GetShopByID(shopID)

	return &Shop{Model: shop}, err
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

func (s *Service) GenerateShopAssetsToken() (string, error) {
	const length = 96
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

	ret := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func (s *Service) SaveNewShopToDB(shop *Shop, u string) error {
	token, err := s.GenerateShopAssetsToken()
	if err != nil {
		return fmt.Errorf("GenerateShopAssetsToken: %w", err)
	}
	shop.Model.AssetsToken = token
	_, err = s.store.SaveNewShopToStore(shop, u)

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
			sr.ID = s.ID
			sr.UserID = s.UserID
			sr.Name = s.ShopName
			sr.AssetsToken = s.AssetsToken

			srList = append(srList, sr)

			srMap[s.ID] = struct{}{}
		}

		created := int64(0)

		if s.PlatformCreated.Valid {
			created = s.PlatformCreated.Time.Unix()
		}

		if s.PlatformID.Valid {
			pMap[s.ID] = append(pMap[s.ID], PlatformReturn{
				ID:      s.PlatformID.String,
				Name:    s.PlatformName.String,
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

func (s *Service) GetProductListJSON(pp []models.Product, paginate *models.Paginate) ([]byte, error) {
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

func (s *Service) GetProductDetailsByID(productID string, shopID string) (*models.ProductDetails, error) {
	pModel, err := s.store.GetProductDetails(productID, shopID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("GetProductByID: %w", err)
	}

	assets, err := s.store.GetProductAssetsByProductID(productID, shopID)
	if err != pgx.ErrNoRows && err != nil {
		return nil, fmt.Errorf("GetProductByID: %w", err)
	}

	mediaList := []models.ProductAssets{}

	for _, a := range assets {
		if a.ID.Valid {
			mediaList = append(mediaList, models.ProductAssets{
				ID:               a.ID.String,
				OriginalFilename: a.OriginalFilename,
				Extension:        a.Extension,
				Status:           a.Status,
				FileSize:         a.FileSize,
				Width:            a.Width,
				Height:           a.Height,
				Order:            a.ImageOrder,
				Created:          a.Created,
			})
		}
	}

	var platformList []models.ProductPlatformInformation
	var platformSKU string

	if pModel.LazadaPlatformSku.Valid {
		platformSKU = strconv.FormatInt(pModel.LazadaPlatformSku.Int64, 10)
	}

	lazadaPlatform := models.ProductPlatformInformation{
		ID:           strconv.FormatInt(pModel.LazadaID.Int64, 10),
		ProductURL:   pModel.LazadaUrl,
		Name:         pModel.LazadaName,
		PlatformName: "lazada",
		PlatformSKU:  platformSKU,
		SellerSKU:    pModel.LazadaSellerSku.String,
		Status:       pModel.LazadaStatus.String,
		SyncStatus:   pModel.LazadaSyncStatus.Bool,
	}

	platformList = append(platformList, lazadaPlatform)

	return &models.ProductDetails{
		Model: &sqlc.Product{
			ID:               pModel.ID,
			Name:             pModel.Name,
			Description:      pModel.Description,
			DescriptionSlate: pModel.DescriptionSlate,
			Msku:             pModel.Msku,
			SellingPrice:     pModel.SellingPrice,
			CostPrice:        pModel.CostPrice,
			ShopID:           pModel.ShopID,
			MediaID:          pModel.MediaID,
			Created:          pModel.Created,
			Updated:          pModel.Updated,
		},
		MediaList: mediaList,
		Platforms: platformList,
	}, nil
}

func (s *Service) GetProductsByNameOrSKU(userID string, name null.String, msku null.String, offset string, limit string) ([]models.Product, models.Paginate, error) {
	var pList []models.Product
	var paginate models.Paginate

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

		mediaList := []string{}
		mediaListString := string(pModel.MediaIDList)
		if mediaListString != "" {
			mediaList = strings.Split(mediaListString, ",")
		}

		var platformList []models.ProductPlatformInformation
		var platformSKU string

		if pModel.LazadaPlatformSku.Valid {
			platformSKU = strconv.FormatInt(pModel.LazadaPlatformSku.Int64, 10)
		}

		lazadaPlatform := models.ProductPlatformInformation{
			ID:           strconv.FormatInt(pModel.LazadaID, 10),
			ProductURL:   pModel.LazadaUrl,
			Name:         pModel.LazadaName,
			PlatformName: "lazada",
			PlatformSKU:  platformSKU,
			SellerSKU:    pModel.LazadaSellerSku,
			Status:       pModel.LazadaStatus.String,
		}

		platformList = append(platformList, lazadaPlatform)

		pList = append(pList, models.Product{
			Model: &sqlc.Product{
				ID:           pModel.ID,
				Name:         pModel.Name,
				Description:  pModel.Description,
				Msku:         pModel.Msku,
				SellingPrice: pModel.SellingPrice,
				CostPrice:    pModel.CostPrice,
				ShopID:       pModel.ShopID,
				MediaID:      pModel.MediaID,
				Created:      pModel.Created,
				Updated:      pModel.Updated,
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

func (s *Service) GetProductsByUserID(userID string, offset string, limit string) ([]models.Product, models.Paginate, error) {
	var pList []models.Product
	var paginate models.Paginate

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

		mediaList := []string{}
		mediaListString := string(pModel.MediaIDList)
		if mediaListString != "" {
			mediaList = strings.Split(mediaListString, ",")
		}

		var platformList []models.ProductPlatformInformation
		var platformSKU string

		if pModel.LazadaPlatformSku.Valid {
			platformSKU = strconv.FormatInt(pModel.LazadaPlatformSku.Int64, 10)
		}

		lazadaPlatform := models.ProductPlatformInformation{
			ID:           strconv.FormatInt(pModel.LazadaID, 10),
			ProductURL:   pModel.LazadaUrl,
			Name:         pModel.LazadaName,
			PlatformName: "lazada",
			PlatformSKU:  platformSKU,
			SellerSKU:    pModel.LazadaSellerSku,
			Status:       pModel.LazadaStatus.String,
		}

		platformList = append(platformList, lazadaPlatform)

		pList = append(pList, models.Product{
			Model: &sqlc.Product{
				ID:           pModel.ID,
				Name:         pModel.Name,
				Description:  pModel.Description,
				Msku:         pModel.Msku,
				SellingPrice: pModel.SellingPrice,
				CostPrice:    pModel.CostPrice,
				ShopID:       pModel.ShopID,
				MediaID:      pModel.MediaID,
				Created:      pModel.Created,
				Updated:      pModel.Updated,
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

func (s *Service) GetSantizedString(d string) string {
	p := bluemonday.StrictPolicy()
	santized := p.Sanitize(d)
	santized = html.UnescapeString(santized)

	return santized
}

func (s *Service) GetLaxoProductFromLazadaData(p *sqlc.ProductsLazada,
	pAttribute *sqlc.ProductsAttributeLazada, pSKU *sqlc.ProductsSkuLazada, shopID string) (*models.Product, error) {

	numericPrice := pgtype.Numeric{}
	numericPrice.Set(pSKU.Price.String)

	sanitzedDescription := s.GetSantizedString(pAttribute.Description.String)

	var slateDescription null.String
	var slateDescriptionString string
	var err error

	if pAttribute.Description.Valid {
		slateDescriptionString, err = s.HTMLToSlate(pAttribute.Description.String, shopID)
		if err != nil && err != ErrEmptyText {
			return nil, fmt.Errorf("HTMLToSlate: %w", err)
		}

		if err == ErrEmptyText {
			slateDescription = null.NewString("", false)
		} else {
			slateDescription = null.StringFrom(slateDescriptionString)
		}
	}

	product := &models.Product{
		Model: &sqlc.Product{
			ID:               "",
			Name:             pAttribute.Name,
			Description:      null.StringFrom(sanitzedDescription),
			DescriptionSlate: slateDescription,
			Msku:             null.StringFrom(pSKU.SellerSku),
			SellingPrice:     numericPrice,
			CostPrice:        pgtype.Numeric{Status: pgtype.Null},
			ShopID:           p.ShopID,
		},
	}

	return product, nil
}

// This should be used when importing products (syncing) from Lazada
func (s *Service) SaveOrUpdateProductToStore(p *models.Product, shopID string, lazadaID string, overwrite bool) (*models.Product, error) {
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
			ProductID:        newModel.ID,
			ProductsLazadaID: null.StringFrom(lazadaID),
		}
		platform, err = s.store.CreateProductPlatform(param)
		if err != nil {
			return nil, fmt.Errorf("CreateProductPlatform: %w", err)
		}

		pReturn = &models.Product{
			Model:         newModel,
			PlatformModel: platform,
		}

		return pReturn, nil
	}

	if !overwrite {
		//@TODO: it's probably better to get the model from DB
		newModel = &sqlc.Product{ID: platform.ProductID}

		pReturn = &models.Product{
			Model:         newModel,
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
		Model:         newModel,
		PlatformModel: platform,
	}

	return pReturn, nil
}
