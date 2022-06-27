package assets

import (
	"encoding/json"

	"laxo.vn/laxo/laxo/sqlc"
)

type AssetRequest struct {
	OriginalName string `json:"originalName"`
	Size         int64  `json:"size"`
	WidthPixels  int64  `json:"width"`
	HeightPixels int64  `json:"height"`
	Hash         string `json:"hash"`
}

type AssignRequest struct {
	Action    string `json:"action"`
	ProductID string `json:"productID"`
	AssetID   string `json:"assetID"`
	Order     int64  `json:"order"`
}

type AssetReply struct {
	Asset  *sqlc.Asset `json:"asset"`
	Upload bool        `json:"upload"`
}

func (a *AssetReply) JSON() ([]byte, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}
