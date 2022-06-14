package assets

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/twmb/murmur3"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store interface {
  GetAssetByMurmur(murmurHash string, shopID string) (*sqlc.Asset, error)
  SaveNewProductMedia(int64, string, []byte, string, string, string) error
  CreateNewAsset(ShopID, MurmurHash, OriginalFilename string, FileSize, Width, Height int64) (*sqlc.Asset, error)
  GetAssetByID(assetID string) (*sqlc.Asset, error)
  SaveAssetToDisk(b []byte, assetID string, shopToken string, fileExt string) (*sqlc.Asset, error)
  GetProductMedia(assetID string, productID string) (*sqlc.ProductsMedia, error)
  CreateProductMedia(assetID string, productID string, status string, order int64) (*sqlc.ProductsMedia, error)
  UpdateProductMedia(assetID string, productID string, status null.String, order null.Int) (*sqlc.ProductsMedia, error)
  UnlinkProductMedia(productID string, assetID string) error
}

var ValidExtenions = map[string]struct{}{
  ".png": {},
  ".jpg": {},
  ".jpeg": {},
  ".jpe": {},
  ".jif": {},
  ".jfif": {},
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server) Service {
  return Service {
    store: store,
    logger: logger,
    server: server,
  }
}

type Service struct {
  store       Store
  logger      *laxo.Logger
  server      *laxo.Server
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

func (s *Service) SaveProductImages(i []string, shopID string, productID string, shopToken string) error {
  for _, v := range i {
    resp, err := http.Get(v)
    if err != nil {
      return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
      return errors.New("lazada response code is not 200")
    }

    b, err := io.ReadAll(resp.Body)
    if err != nil {
      return err
    }

    h64 := murmur3.New64()
    h64.Write(b)

    mID := int64(h64.Sum64())

    //_, err = s.GetProductMediaByMurmur(mID, productID)
    //if err != pgx.ErrNoRows && err != nil {
    //  return err
    //}

    if err == pgx.ErrNoRows {
      filename := path.Base(v)
      err = s.store.SaveNewProductMedia(mID, filename, b, shopID, productID, shopToken)
      if err != nil {
        return err
      }
    }
  }
  return nil
}
