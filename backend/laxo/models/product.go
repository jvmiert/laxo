package models

import (
	"encoding/json"

	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/sqlc"
)

type Element struct {
	Type      string    `json:"type,omitempty"`
	Src       string    `json:"src,omitempty"`
	Width     int64     `json:"width,omitempty"`
	Height    int64     `json:"height,omitempty"`
	Align     string    `json:"align,omitempty"`
	Text      string    `json:"text,omitempty"`
	Bold      bool      `json:"bold,omitempty"`
	Underline bool      `json:"underline,omitempty"`
	Italic    bool      `json:"italic,omitempty"`
	Children  []Element `json:"children,omitempty"`
}

type ProductPlatformInformation struct {
	ID           string      `json:"id"`
	SellerSKU    string      `json:"sellerSKU"`
	PlatformSKU  string      `json:"platformSKU"`
	PlatformName string      `json:"platformName"`
	Status       string      `json:"status"`
	SyncStatus   bool        `json:"syncStatus"`
	Name         null.String `json:"name"`
	ProductURL   null.String `json:"productURL"`
}

type ProductChangedSyncRequest struct {
	Platform string `json:"platform"`
	State    bool   `json:"state"`
}

type ProductDetailPostRequest struct {
	SellingPrice int         `json:"sellingPrice"`
	CostPrice    int         `json:"costPrice"`
	Name         null.String `json:"name"`
	Description  []Element   `json:"description"`
	Msku         string      `json:"msku"`
}

type Product struct {
	Model         *sqlc.Product                `json:"product"`
	MediaModels   []sqlc.ProductsMedia         `json:"-"`
	PlatformModel *sqlc.ProductsPlatform       `json:"-"`
	MediaList     []string                     `json:"mediaList"`
	Platforms     []ProductPlatformInformation `json:"platforms"`
}

func (p *Product) JSON() ([]byte, error) {
	bytes, err := json.Marshal(p)

	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

type ProductAssets struct {
	ID               string      `json:"id"`
	OriginalFilename null.String `json:"originalFilename"`
	Extension        null.String `json:"extension"`
	Status           null.String `json:"status"`
	FileSize         null.Int    `json:"fileSize"`
	Order            null.Int    `json:"order"`
	Width            null.Int    `json:"width"`
	Height           null.Int    `json:"height"`
}

type ProductDetails struct {
	Model         *sqlc.Product                `json:"product"`
	MediaModels   []sqlc.ProductsMedia         `json:"-"`
	PlatformModel *sqlc.ProductsPlatform       `json:"-"`
	MediaList     []ProductAssets              `json:"mediaList"`
	Platforms     []ProductPlatformInformation `json:"platforms"`
}

func (p *ProductDetails) JSON() ([]byte, error) {
	bytes, err := json.Marshal(p)

	if err != nil {
		return bytes, err
	}

	return bytes, nil
}
