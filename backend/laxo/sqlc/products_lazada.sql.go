// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: products_lazada.sql

package sqlc

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

const createLazadaProduct = `-- name: CreateLazadaProduct :one
INSERT INTO products_lazada (
  lazada_id, lazada_primary_category, created, updated,
  status, sub_status, shop_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, lazada_id, lazada_primary_category, created, updated, status, sub_status, shop_id
`

type CreateLazadaProductParams struct {
	LazadaID              int64       `json:"lazadaID"`
	LazadaPrimaryCategory int64       `json:"lazadaPrimaryCategory"`
	Created               time.Time   `json:"created"`
	Updated               time.Time   `json:"updated"`
	Status                null.String `json:"status"`
	SubStatus             null.String `json:"subStatus"`
	ShopID                string      `json:"shopID"`
}

func (q *Queries) CreateLazadaProduct(ctx context.Context, arg CreateLazadaProductParams) (ProductsLazada, error) {
	row := q.db.QueryRow(ctx, createLazadaProduct,
		arg.LazadaID,
		arg.LazadaPrimaryCategory,
		arg.Created,
		arg.Updated,
		arg.Status,
		arg.SubStatus,
		arg.ShopID,
	)
	var i ProductsLazada
	err := row.Scan(
		&i.ID,
		&i.LazadaID,
		&i.LazadaPrimaryCategory,
		&i.Created,
		&i.Updated,
		&i.Status,
		&i.SubStatus,
		&i.ShopID,
	)
	return i, err
}

const createLazadaProductAttribute = `-- name: CreateLazadaProductAttribute :one
INSERT INTO products_attribute_lazada (
  name, short_description, description, brand, model, headphone_features,
  bluetooth, warranty_type, warranty, hazmat, expire_date,
  brand_classification, ingredient_preference, lot_number, units_hb,
  fmlt_skincare, quantitative, skincare_by_age, skin_benefit, skin_type,
  user_manual, country_origin_hb, color_family, fragrance_family,
  source, product_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9, $10, $11, $12, $13, $14,
  $15, $16, $17, $18, $19, $20, $21,
  $22, $23, $24, $25, $26
)
RETURNING id, name, short_description, description, brand, model, headphone_features, bluetooth, warranty_type, warranty, hazmat, expire_date, brand_classification, ingredient_preference, lot_number, units_hb, fmlt_skincare, quantitative, skincare_by_age, skin_benefit, skin_type, user_manual, country_origin_hb, color_family, fragrance_family, source, product_id
`

type CreateLazadaProductAttributeParams struct {
	Name                 null.String `json:"name"`
	ShortDescription     null.String `json:"shortDescription"`
	Description          null.String `json:"description"`
	Brand                null.String `json:"brand"`
	Model                null.String `json:"model"`
	HeadphoneFeatures    null.String `json:"headphoneFeatures"`
	Bluetooth            null.String `json:"bluetooth"`
	WarrantyType         null.String `json:"warrantyType"`
	Warranty             null.String `json:"warranty"`
	Hazmat               null.String `json:"hazmat"`
	ExpireDate           null.String `json:"expireDate"`
	BrandClassification  null.String `json:"brandClassification"`
	IngredientPreference null.String `json:"ingredientPreference"`
	LotNumber            null.String `json:"lotNumber"`
	UnitsHb              null.String `json:"unitsHb"`
	FmltSkincare         null.String `json:"fmltSkincare"`
	Quantitative         null.String `json:"quantitative"`
	SkincareByAge        null.String `json:"skincareByAge"`
	SkinBenefit          null.String `json:"skinBenefit"`
	SkinType             null.String `json:"skinType"`
	UserManual           null.String `json:"userManual"`
	CountryOriginHb      null.String `json:"countryOriginHb"`
	ColorFamily          null.String `json:"colorFamily"`
	FragranceFamily      null.String `json:"fragranceFamily"`
	Source               null.String `json:"source"`
	ProductID            string      `json:"productID"`
}

func (q *Queries) CreateLazadaProductAttribute(ctx context.Context, arg CreateLazadaProductAttributeParams) (ProductsAttributeLazada, error) {
	row := q.db.QueryRow(ctx, createLazadaProductAttribute,
		arg.Name,
		arg.ShortDescription,
		arg.Description,
		arg.Brand,
		arg.Model,
		arg.HeadphoneFeatures,
		arg.Bluetooth,
		arg.WarrantyType,
		arg.Warranty,
		arg.Hazmat,
		arg.ExpireDate,
		arg.BrandClassification,
		arg.IngredientPreference,
		arg.LotNumber,
		arg.UnitsHb,
		arg.FmltSkincare,
		arg.Quantitative,
		arg.SkincareByAge,
		arg.SkinBenefit,
		arg.SkinType,
		arg.UserManual,
		arg.CountryOriginHb,
		arg.ColorFamily,
		arg.FragranceFamily,
		arg.Source,
		arg.ProductID,
	)
	var i ProductsAttributeLazada
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ShortDescription,
		&i.Description,
		&i.Brand,
		&i.Model,
		&i.HeadphoneFeatures,
		&i.Bluetooth,
		&i.WarrantyType,
		&i.Warranty,
		&i.Hazmat,
		&i.ExpireDate,
		&i.BrandClassification,
		&i.IngredientPreference,
		&i.LotNumber,
		&i.UnitsHb,
		&i.FmltSkincare,
		&i.Quantitative,
		&i.SkincareByAge,
		&i.SkinBenefit,
		&i.SkinType,
		&i.UserManual,
		&i.CountryOriginHb,
		&i.ColorFamily,
		&i.FragranceFamily,
		&i.Source,
		&i.ProductID,
	)
	return i, err
}

const createLazadaProductSKU = `-- name: CreateLazadaProductSKU :one
INSERT INTO products_sku_lazada (
  status, quantity, seller_sku, shop_sku, sku_id, url,
  price, available, package_content, package_width, package_weight,
  package_length, package_height, special_price, special_to_time,
  special_from_time, special_from_date, special_to_date, product_id, shop_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9, $10, $11, $12, $13, $14,
  $15, $16, $17, $18, $19, $20
)
RETURNING id, status, quantity, seller_sku, shop_sku, sku_id, url, price, available, package_content, package_width, package_weight, package_length, package_height, special_price, special_to_time, special_from_time, special_from_date, special_to_date, product_id, shop_id
`

type CreateLazadaProductSKUParams struct {
	Status          null.String `json:"status"`
	Quantity        null.Int    `json:"quantity"`
	SellerSku       string      `json:"sellerSku"`
	ShopSku         string      `json:"shopSku"`
	SkuID           null.Int    `json:"skuID"`
	Url             null.String `json:"url"`
	Price           null.String `json:"price"`
	Available       null.Int    `json:"available"`
	PackageContent  null.String `json:"packageContent"`
	PackageWidth    null.String `json:"packageWidth"`
	PackageWeight   null.String `json:"packageWeight"`
	PackageLength   null.String `json:"packageLength"`
	PackageHeight   null.String `json:"packageHeight"`
	SpecialPrice    null.String `json:"specialPrice"`
	SpecialToTime   null.Time   `json:"specialToTime"`
	SpecialFromTime null.Time   `json:"specialFromTime"`
	SpecialFromDate null.Time   `json:"specialFromDate"`
	SpecialToDate   null.Time   `json:"specialToDate"`
	ProductID       string      `json:"productID"`
	ShopID          string      `json:"shopID"`
}

func (q *Queries) CreateLazadaProductSKU(ctx context.Context, arg CreateLazadaProductSKUParams) (ProductsSkuLazada, error) {
	row := q.db.QueryRow(ctx, createLazadaProductSKU,
		arg.Status,
		arg.Quantity,
		arg.SellerSku,
		arg.ShopSku,
		arg.SkuID,
		arg.Url,
		arg.Price,
		arg.Available,
		arg.PackageContent,
		arg.PackageWidth,
		arg.PackageWeight,
		arg.PackageLength,
		arg.PackageHeight,
		arg.SpecialPrice,
		arg.SpecialToTime,
		arg.SpecialFromTime,
		arg.SpecialFromDate,
		arg.SpecialToDate,
		arg.ProductID,
		arg.ShopID,
	)
	var i ProductsSkuLazada
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Quantity,
		&i.SellerSku,
		&i.ShopSku,
		&i.SkuID,
		&i.Url,
		&i.Price,
		&i.Available,
		&i.PackageContent,
		&i.PackageWidth,
		&i.PackageWeight,
		&i.PackageLength,
		&i.PackageHeight,
		&i.SpecialPrice,
		&i.SpecialToTime,
		&i.SpecialFromTime,
		&i.SpecialFromDate,
		&i.SpecialToDate,
		&i.ProductID,
		&i.ShopID,
	)
	return i, err
}

const getLazadaProductAttributeByProductID = `-- name: GetLazadaProductAttributeByProductID :one
SELECT id, name, short_description, description, brand, model, headphone_features, bluetooth, warranty_type, warranty, hazmat, expire_date, brand_classification, ingredient_preference, lot_number, units_hb, fmlt_skincare, quantitative, skincare_by_age, skin_benefit, skin_type, user_manual, country_origin_hb, color_family, fragrance_family, source, product_id FROM products_attribute_lazada
WHERE product_id = $1
LIMIT 1
`

func (q *Queries) GetLazadaProductAttributeByProductID(ctx context.Context, productID string) (ProductsAttributeLazada, error) {
	row := q.db.QueryRow(ctx, getLazadaProductAttributeByProductID, productID)
	var i ProductsAttributeLazada
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ShortDescription,
		&i.Description,
		&i.Brand,
		&i.Model,
		&i.HeadphoneFeatures,
		&i.Bluetooth,
		&i.WarrantyType,
		&i.Warranty,
		&i.Hazmat,
		&i.ExpireDate,
		&i.BrandClassification,
		&i.IngredientPreference,
		&i.LotNumber,
		&i.UnitsHb,
		&i.FmltSkincare,
		&i.Quantitative,
		&i.SkincareByAge,
		&i.SkinBenefit,
		&i.SkinType,
		&i.UserManual,
		&i.CountryOriginHb,
		&i.ColorFamily,
		&i.FragranceFamily,
		&i.Source,
		&i.ProductID,
	)
	return i, err
}

const getLazadaProductByLazadaID = `-- name: GetLazadaProductByLazadaID :one
SELECT id, lazada_id, lazada_primary_category, created, updated, status, sub_status, shop_id FROM products_lazada
WHERE lazada_id = $1 AND shop_id = $2
LIMIT 1
`

type GetLazadaProductByLazadaIDParams struct {
	LazadaID int64  `json:"lazadaID"`
	ShopID   string `json:"shopID"`
}

func (q *Queries) GetLazadaProductByLazadaID(ctx context.Context, arg GetLazadaProductByLazadaIDParams) (ProductsLazada, error) {
	row := q.db.QueryRow(ctx, getLazadaProductByLazadaID, arg.LazadaID, arg.ShopID)
	var i ProductsLazada
	err := row.Scan(
		&i.ID,
		&i.LazadaID,
		&i.LazadaPrimaryCategory,
		&i.Created,
		&i.Updated,
		&i.Status,
		&i.SubStatus,
		&i.ShopID,
	)
	return i, err
}

const getLazadaProductSKUByProductID = `-- name: GetLazadaProductSKUByProductID :one
SELECT id, status, quantity, seller_sku, shop_sku, sku_id, url, price, available, package_content, package_width, package_weight, package_length, package_height, special_price, special_to_time, special_from_time, special_from_date, special_to_date, product_id, shop_id FROM products_sku_lazada
WHERE product_id = $1
LIMIT 1
`

func (q *Queries) GetLazadaProductSKUByProductID(ctx context.Context, productID string) (ProductsSkuLazada, error) {
	row := q.db.QueryRow(ctx, getLazadaProductSKUByProductID, productID)
	var i ProductsSkuLazada
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Quantity,
		&i.SellerSku,
		&i.ShopSku,
		&i.SkuID,
		&i.Url,
		&i.Price,
		&i.Available,
		&i.PackageContent,
		&i.PackageWidth,
		&i.PackageWeight,
		&i.PackageLength,
		&i.PackageHeight,
		&i.SpecialPrice,
		&i.SpecialToTime,
		&i.SpecialFromTime,
		&i.SpecialFromDate,
		&i.SpecialToDate,
		&i.ProductID,
		&i.ShopID,
	)
	return i, err
}

const updateLazadaProduct = `-- name: UpdateLazadaProduct :one
UPDATE products_lazada SET
  lazada_id = $1, lazada_primary_category = $2, created = $3, updated = $4,
  status = $5, sub_status = $6
WHERE id = $7
RETURNING id, lazada_id, lazada_primary_category, created, updated, status, sub_status, shop_id
`

type UpdateLazadaProductParams struct {
	LazadaID              int64       `json:"lazadaID"`
	LazadaPrimaryCategory int64       `json:"lazadaPrimaryCategory"`
	Created               time.Time   `json:"created"`
	Updated               time.Time   `json:"updated"`
	Status                null.String `json:"status"`
	SubStatus             null.String `json:"subStatus"`
	ID                    string      `json:"id"`
}

func (q *Queries) UpdateLazadaProduct(ctx context.Context, arg UpdateLazadaProductParams) (ProductsLazada, error) {
	row := q.db.QueryRow(ctx, updateLazadaProduct,
		arg.LazadaID,
		arg.LazadaPrimaryCategory,
		arg.Created,
		arg.Updated,
		arg.Status,
		arg.SubStatus,
		arg.ID,
	)
	var i ProductsLazada
	err := row.Scan(
		&i.ID,
		&i.LazadaID,
		&i.LazadaPrimaryCategory,
		&i.Created,
		&i.Updated,
		&i.Status,
		&i.SubStatus,
		&i.ShopID,
	)
	return i, err
}

const updateLazadaProductAttribute = `-- name: UpdateLazadaProductAttribute :one
UPDATE products_attribute_lazada SET
  name = $1, short_description = $2, description = $3, brand = $4, model = $5,
  headphone_features = $6, bluetooth = $7, warranty_type = $8, warranty = $9,
  hazmat = $10, expire_date = $11, brand_classification = $12,
  ingredient_preference = $13, lot_number = $14, units_hb = $15,
  fmlt_skincare = $16, quantitative = $17, skincare_by_age = $18,
  skin_benefit = $19, skin_type = $20, user_manual = $21,
  country_origin_hb = $22, color_family = $23, fragrance_family = $24,
  source = $25
WHERE id = $26
RETURNING id, name, short_description, description, brand, model, headphone_features, bluetooth, warranty_type, warranty, hazmat, expire_date, brand_classification, ingredient_preference, lot_number, units_hb, fmlt_skincare, quantitative, skincare_by_age, skin_benefit, skin_type, user_manual, country_origin_hb, color_family, fragrance_family, source, product_id
`

type UpdateLazadaProductAttributeParams struct {
	Name                 null.String `json:"name"`
	ShortDescription     null.String `json:"shortDescription"`
	Description          null.String `json:"description"`
	Brand                null.String `json:"brand"`
	Model                null.String `json:"model"`
	HeadphoneFeatures    null.String `json:"headphoneFeatures"`
	Bluetooth            null.String `json:"bluetooth"`
	WarrantyType         null.String `json:"warrantyType"`
	Warranty             null.String `json:"warranty"`
	Hazmat               null.String `json:"hazmat"`
	ExpireDate           null.String `json:"expireDate"`
	BrandClassification  null.String `json:"brandClassification"`
	IngredientPreference null.String `json:"ingredientPreference"`
	LotNumber            null.String `json:"lotNumber"`
	UnitsHb              null.String `json:"unitsHb"`
	FmltSkincare         null.String `json:"fmltSkincare"`
	Quantitative         null.String `json:"quantitative"`
	SkincareByAge        null.String `json:"skincareByAge"`
	SkinBenefit          null.String `json:"skinBenefit"`
	SkinType             null.String `json:"skinType"`
	UserManual           null.String `json:"userManual"`
	CountryOriginHb      null.String `json:"countryOriginHb"`
	ColorFamily          null.String `json:"colorFamily"`
	FragranceFamily      null.String `json:"fragranceFamily"`
	Source               null.String `json:"source"`
	ID                   string      `json:"id"`
}

func (q *Queries) UpdateLazadaProductAttribute(ctx context.Context, arg UpdateLazadaProductAttributeParams) (ProductsAttributeLazada, error) {
	row := q.db.QueryRow(ctx, updateLazadaProductAttribute,
		arg.Name,
		arg.ShortDescription,
		arg.Description,
		arg.Brand,
		arg.Model,
		arg.HeadphoneFeatures,
		arg.Bluetooth,
		arg.WarrantyType,
		arg.Warranty,
		arg.Hazmat,
		arg.ExpireDate,
		arg.BrandClassification,
		arg.IngredientPreference,
		arg.LotNumber,
		arg.UnitsHb,
		arg.FmltSkincare,
		arg.Quantitative,
		arg.SkincareByAge,
		arg.SkinBenefit,
		arg.SkinType,
		arg.UserManual,
		arg.CountryOriginHb,
		arg.ColorFamily,
		arg.FragranceFamily,
		arg.Source,
		arg.ID,
	)
	var i ProductsAttributeLazada
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ShortDescription,
		&i.Description,
		&i.Brand,
		&i.Model,
		&i.HeadphoneFeatures,
		&i.Bluetooth,
		&i.WarrantyType,
		&i.Warranty,
		&i.Hazmat,
		&i.ExpireDate,
		&i.BrandClassification,
		&i.IngredientPreference,
		&i.LotNumber,
		&i.UnitsHb,
		&i.FmltSkincare,
		&i.Quantitative,
		&i.SkincareByAge,
		&i.SkinBenefit,
		&i.SkinType,
		&i.UserManual,
		&i.CountryOriginHb,
		&i.ColorFamily,
		&i.FragranceFamily,
		&i.Source,
		&i.ProductID,
	)
	return i, err
}

const updateLazadaProductSKU = `-- name: UpdateLazadaProductSKU :one
UPDATE products_sku_lazada SET
  status = $1, quantity = $2, seller_sku = $3, shop_sku = $4, sku_id = $5, url = $6,
  price = $7, available = $8, package_content = $9, package_width = $10, package_weight = $11,
  package_length = $12, package_height = $13, special_price = $14, special_to_time = $15,
  special_from_time = $16, special_from_date = $17, special_to_date = $18, product_id = $19, shop_id = $20
WHERE id = $21
RETURNING id, status, quantity, seller_sku, shop_sku, sku_id, url, price, available, package_content, package_width, package_weight, package_length, package_height, special_price, special_to_time, special_from_time, special_from_date, special_to_date, product_id, shop_id
`

type UpdateLazadaProductSKUParams struct {
	Status          null.String `json:"status"`
	Quantity        null.Int    `json:"quantity"`
	SellerSku       string      `json:"sellerSku"`
	ShopSku         string      `json:"shopSku"`
	SkuID           null.Int    `json:"skuID"`
	Url             null.String `json:"url"`
	Price           null.String `json:"price"`
	Available       null.Int    `json:"available"`
	PackageContent  null.String `json:"packageContent"`
	PackageWidth    null.String `json:"packageWidth"`
	PackageWeight   null.String `json:"packageWeight"`
	PackageLength   null.String `json:"packageLength"`
	PackageHeight   null.String `json:"packageHeight"`
	SpecialPrice    null.String `json:"specialPrice"`
	SpecialToTime   null.Time   `json:"specialToTime"`
	SpecialFromTime null.Time   `json:"specialFromTime"`
	SpecialFromDate null.Time   `json:"specialFromDate"`
	SpecialToDate   null.Time   `json:"specialToDate"`
	ProductID       string      `json:"productID"`
	ShopID          string      `json:"shopID"`
	ID              string      `json:"id"`
}

func (q *Queries) UpdateLazadaProductSKU(ctx context.Context, arg UpdateLazadaProductSKUParams) (ProductsSkuLazada, error) {
	row := q.db.QueryRow(ctx, updateLazadaProductSKU,
		arg.Status,
		arg.Quantity,
		arg.SellerSku,
		arg.ShopSku,
		arg.SkuID,
		arg.Url,
		arg.Price,
		arg.Available,
		arg.PackageContent,
		arg.PackageWidth,
		arg.PackageWeight,
		arg.PackageLength,
		arg.PackageHeight,
		arg.SpecialPrice,
		arg.SpecialToTime,
		arg.SpecialFromTime,
		arg.SpecialFromDate,
		arg.SpecialToDate,
		arg.ProductID,
		arg.ShopID,
		arg.ID,
	)
	var i ProductsSkuLazada
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Quantity,
		&i.SellerSku,
		&i.ShopSku,
		&i.SkuID,
		&i.Url,
		&i.Price,
		&i.Available,
		&i.PackageContent,
		&i.PackageWidth,
		&i.PackageWeight,
		&i.PackageLength,
		&i.PackageHeight,
		&i.SpecialPrice,
		&i.SpecialToTime,
		&i.SpecialFromTime,
		&i.SpecialFromDate,
		&i.SpecialToDate,
		&i.ProductID,
		&i.ShopID,
	)
	return i, err
}
