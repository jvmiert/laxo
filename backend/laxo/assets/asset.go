package assets

import (
	"encoding/json"

	"laxo.vn/laxo/laxo/sqlc"
)

type AssetRequest struct {
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	WidthPixels  int64  `json:"width"`
	HeightPixels int64  `json:"height"`
	Hash         string `json:"hash"`
}

type AssetReply struct {
	Asset     *sqlc.Asset `json:"asset"`
	Upload    bool        `json:"upload"`
}

func (a *AssetReply) JSON() ([]byte, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}
