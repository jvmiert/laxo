package assets

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/jackc/pgx/v4"
	"github.com/twmb/murmur3"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store interface {
  GetAssetByMurmur(murmurHash string, shopID string) (*sqlc.Asset, error)
  SaveNewProductMedia(int64, string, []byte, string, string, string) error
  CreateNewAsset(ShopID, MurmurHash, OriginalFilename string, FileSize, Width, Height int64) (*sqlc.Asset, error)
  GetAssetByID(assetID string) (*sqlc.Asset, error)
  SaveAssetToDisk(b []byte, assetID string, shopToken string) (*sqlc.Asset, error)
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

func (s *Service) SaveAssetToDisk(b []byte, assetID string, shopToken string) (*sqlc.Asset, error) {
  //@TODO: Move the extension checking out of the store and into here so we can validate the file extension
  return s.store.SaveAssetToDisk(b, assetID, shopToken)
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
