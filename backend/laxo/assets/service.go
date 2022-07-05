package assets

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/twmb/murmur3"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/models"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store interface {
	GetAssetByMurmur(murmurHash string, shopID string) (*sqlc.Asset, error)
	CreateNewAsset(ShopID, MurmurHash, OriginalFilename string, FileSize, Width, Height int64) (*sqlc.Asset, error)
	GetAssetByID(assetID string) (*sqlc.Asset, error)
	SaveAssetToDisk(b []byte, assetID string, shopToken string, fileExt string) (*sqlc.Asset, error)
	GetProductMedia(assetID string, productID string) (*sqlc.ProductsMedia, error)
	CreateProductMedia(assetID string, productID string, status string, order int64) (*sqlc.ProductsMedia, error)
	UpdateProductMedia(assetID string, productID string, status null.String, order null.Int) (*sqlc.ProductsMedia, error)
	UnlinkProductMedia(productID string, assetID string) error
	GetAllAssetsByShopID(shopID string, limit int32, offset int32) ([]sqlc.GetAllAssetsByShopIDRow, error)
	GetAssetBytesByID(assetID string, shopID string, shopToken string) ([]byte, error)
	GetAssetRankByIDAndShopID(assetID string, shopID string) (int64, error)
}

var ErrImageURLForbidden = errors.New("image url returned forbidden")

var ValidExtenions = map[string]struct{}{
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".jpe":  {},
	".jif":  {},
	".jfif": {},
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server) Service {
	return Service{
		store:  store,
		logger: logger,
		server: server,
	}
}

type Service struct {
	store  Store
	logger *laxo.Logger
	server *laxo.Server
}

func (s *Service) GetAssetRankByIDAndShopID(assetID string, shopID string) ([]byte, error) {
	rank, err := s.store.GetAssetRankByIDAndShopID(assetID, shopID)
	if err != nil {
		return nil, fmt.Errorf("GetAssetRankByIDAndShopID: %w", err)
	}

	returnObject := map[string]interface{}{
		"rank": rank,
	}

	bytes, err := json.Marshal(returnObject)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

func (s *Service) GetAssetBytesByID(assetID string, shopID string, shopToken string) ([]byte, error) {
	b, err := s.store.GetAssetBytesByID(assetID, shopID, shopToken)
	if err != nil {
		return nil, fmt.Errorf("GetAssetBytesByID: %w", err)
	}

	return b, err
}

func (s *Service) ExtractImagesFromDescription(d string, shopID string, assetsToken string) error {
	findImages := regexp.MustCompile(`src=["'](.*?)["']`)

	matches := findImages.FindAllStringSubmatch(d, -1)

	for _, element := range matches {
		if len(element) == 0 {
			continue
		}

		url := element[1]
		b, err := s.ExtractImageFromURL(url)
		if err != nil {
			if errors.Is(err, ErrImageURLForbidden) {
				s.server.Logger.Errorw("skipping image due to forbidden return", "url", url)
				continue
			}
			return fmt.Errorf("ExtractImageFromURL: %w", err)
		}

		w, h, err := s.GetImageWidthHeight(b)
		if err != nil {
			return fmt.Errorf("GetImageWidthHeight: %w", err)
		}

		mID := s.GetMurmurFromBytes(b)

		r, err := s.GetOrCreateAsset(AssetRequest{
			OriginalName: path.Base(url),
			Size:         int64(len(b)),
			WidthPixels:  int64(w),
			HeightPixels: int64(h),
			Hash:         mID,
		}, shopID)
		if err != nil {
			return fmt.Errorf("GetOrCreateAsset: %w", err)
		}

		if r.Upload {
			_, err = s.SaveAssetToDisk(b, r.Asset.ID, assetsToken)
			if err != nil {
				return fmt.Errorf("SaveAssetToDisk: %w", err)
			}
		}
	}

	return nil
}

func (s *Service) ValidAssignReply() ([]byte, error) {
	reply := make(map[string]bool)
	reply["error"] = false

	return json.Marshal(reply)
}

func (s *Service) ActivateAssetAssignment(r *AssignRequest) error {
	_, err := s.store.GetProductMedia(r.AssetID, r.ProductID)
	if err != pgx.ErrNoRows && err != nil {
		return fmt.Errorf("GetProductMedia: %w", err)
	}

	if err == pgx.ErrNoRows {
		_, err := s.store.CreateProductMedia(r.AssetID, r.ProductID, "active", r.Order)
		if err != nil {
			return fmt.Errorf("CreateProductMedia: %w", err)
		}
	} else {
		_, err := s.store.UpdateProductMedia(r.AssetID, r.ProductID, null.StringFrom("active"), null.IntFrom(r.Order))
		if err != nil {
			return fmt.Errorf("UpdateProductMedia: %w", err)
		}
	}

	return nil
}

func (s *Service) UnlinkProductMedia(r *AssignRequest) error {
	err := s.store.UnlinkProductMedia(r.ProductID, r.AssetID)
	if err != nil {
		return fmt.Errorf("UnlinkProductMedia: %w", err)
	}
	return nil
}

func (s *Service) ModifyAssetAssignment(r *AssignRequest) error {
	if r.Action == "active" {
		err := s.ActivateAssetAssignment(r)
		if err != nil {
			return fmt.Errorf("ActivateAssetAssignment: %w", err)
		}
		return nil
	}

	if r.Action == "inactive" {
		//@TODO: implement
	}

	if r.Action == "delete" {
		err := s.UnlinkProductMedia(r)
		if err != nil {
			return fmt.Errorf("UnLinkImageAsset: %w", err)
		}
		return nil
	}

	return nil
}

func (s *Service) ValidateAssetExtension(fileName string) error {
	extension := strings.ToLower(filepath.Ext(fileName))

	if _, ok := ValidExtenions[extension]; !ok {
		return errors.New("invalid extension")
	}
	return nil
}

func (s *Service) SaveAssetToDisk(b []byte, assetID string, shopToken string) (*sqlc.Asset, error) {
	filetype := http.DetectContentType(b)

	ext, err := mime.ExtensionsByType(filetype)
	if err != nil {
		return nil, err
	}

	if len(ext) == 0 {
		return nil, errors.New("no extension found for image")
	}

	fileExt := ""

	for _, e := range ext {
		if e == ".jpg" {
			fileExt = ".jpg"
			break
		}
	}

	if fileExt == "" {
		fileExt = ext[0]
	}

	return s.store.SaveAssetToDisk(b, assetID, shopToken, fileExt)
}

func (s *Service) AssetJSON(a *sqlc.Asset) ([]byte, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

func (s *Service) ValidateAssetHash(hex string, assetID string) error {
	a, err := s.store.GetAssetByID(assetID)
	if err != nil {
		return fmt.Errorf("GetAssetByID: %w", err)
	}

	if a.MurmurHash != hex {
		return errors.New("asset hash does not match")
	}

	return nil
}

func (s *Service) GetMurmurFromBytes(b []byte) string {
	h64 := murmur3.New128()
	h64.Write(b)

	mID1, mID2 := h64.Sum128()

	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], mID1)
	binary.BigEndian.PutUint64(buf[8:], mID2)
	hex := hex.EncodeToString(buf)

	return hex
}

func (s *Service) ExtractImageFromRequest(w http.ResponseWriter, b io.ReadCloser) ([]byte, error) {
	newB := http.MaxBytesReader(w, b, int64(s.server.Config.MaxAssetSize))
	return ioutil.ReadAll(newB)
}

func (s *Service) ExtractImageFromURL(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 403 {
			return nil, ErrImageURLForbidden
		}
		s.server.Logger.Errorw("get response code != 200", "status", resp.Status, "url", u)
		return nil, errors.New("get response code is not 200")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) GetOrCreateAsset(request AssetRequest, shopID string) (*AssetReply, error) {
	var a AssetReply

	asset, err := s.store.GetAssetByMurmur(request.Hash, shopID)
	if err != pgx.ErrNoRows && err != nil {
		return nil, fmt.Errorf("GetAssetByMurmur: %w", err)
	}

	if err == pgx.ErrNoRows {
		assetModel, err := s.store.CreateNewAsset(
			shopID, request.Hash, request.OriginalName,
			request.Size, request.WidthPixels, request.HeightPixels,
		)
		if err != nil {
			return nil, fmt.Errorf("CreateNewAsset: %w", err)
		}

		a.Asset = assetModel
		a.Upload = true
		return &a, nil
	}

	a.Asset = asset
	a.Upload = false
	return &a, nil
}

func (s *Service) DecodeImage(b []byte) (*image.Config, error) {
	r := bytes.NewReader(b)

	img, _, err := image.DecodeConfig(r)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func (s *Service) GetImageWidthHeight(b []byte) (width int, height int, err error) {
	imgConfig, err := s.DecodeImage(b)
	if err != nil {
		return 0, 0, fmt.Errorf("DecodeImage: %w", err)
	}

	return imgConfig.Width, imgConfig.Height, nil
}

func (s *Service) SaveAndAssignImagesFromURL(i []string, shopID string, productID string, assetsToken string) error {
	for _, v := range i {
		b, err := s.ExtractImageFromURL(v)
		if err != nil {
			if errors.Is(err, ErrImageURLForbidden) {
				s.server.Logger.Errorw("skipping image due to forbidden return", "url", v)
				continue
			}
			return fmt.Errorf("ExtractImageFromURL: %w", err)
		}

		w, h, err := s.GetImageWidthHeight(b)
		if err != nil {
			return fmt.Errorf("GetImageWidthHeight: %w", err)
		}

		mID := s.GetMurmurFromBytes(b)

		r, err := s.GetOrCreateAsset(AssetRequest{
			OriginalName: path.Base(v),
			Size:         int64(len(b)),
			WidthPixels:  int64(w),
			HeightPixels: int64(h),
			Hash:         mID,
		}, shopID)
		if err != nil {
			return fmt.Errorf("GetOrCreateAsset: %w", err)
		}

		if r.Upload {
			_, err = s.SaveAssetToDisk(b, r.Asset.ID, assetsToken)
			if err != nil {
				return fmt.Errorf("SaveAssetToDisk: %w", err)
			}
		}

		err = s.ActivateAssetAssignment(&AssignRequest{
			Action:    "active",
			ProductID: productID,
			AssetID:   r.Asset.ID,
			Order:     0,
		})
		if err != nil {
			return fmt.Errorf("ActivateAssetAssignment: %w", err)
		}
	}

	return nil
}

func (s *Service) GetAllAssetsByShopID(shopID string, offset string, limit string) ([]sqlc.Asset, models.Paginate, error) {
	var aList []sqlc.Asset
	var paginate models.Paginate

	offsetI, err := strconv.Atoi(offset)
	if err != nil {
		return aList, paginate, fmt.Errorf("atoi offset: %w", err)
	}

	limitI, err := strconv.Atoi(limit)
	if err != nil {
		return aList, paginate, fmt.Errorf("atoi limit: %w", err)
	}

	rows, err := s.store.GetAllAssetsByShopID(shopID, int32(limitI), int32(offsetI))
	if err != nil {
		return aList, paginate, fmt.Errorf("GetProductsByShopID: %w", err)
	}

	total := int64(0)
	if len(rows) > 0 {
		total = rows[0].Count
	}

	for _, row := range rows {
		asset := sqlc.Asset{
			ID:               row.ID,
			ShopID:           row.ShopID,
			MurmurHash:       row.MurmurHash,
			OriginalFilename: row.OriginalFilename,
			Extension:        row.Extension,
			FileSize:         row.FileSize,
			Width:            row.Width,
			Height:           row.Height,
			Created:          row.Created,
		}

		aList = append(aList, asset)
	}

	paginate.Total = total
	paginate.Pages = (total + int64(limitI) - 1) / int64(limitI)
	paginate.Limit = int64(limitI)
	paginate.Offset = int64(offsetI)

	return aList, paginate, nil
}

func (s *Service) GetAssetJSON(pp []sqlc.Asset, paginate *models.Paginate) ([]byte, error) {
	aList := []json.RawMessage{}

	for _, p := range pp {
		b, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		j := json.RawMessage(b)
		aList = append(aList, j)
	}

	productData := map[string]interface{}{
		"assets":   aList,
		"paginate": paginate,
	}

	bytes, err := json.Marshal(productData)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}
