// Code generated by sqlc. DO NOT EDIT.

package sqlc

import (
	"database/sql"
	"time"
)

type PlatformLazada struct {
	ID               string       `json:"id"`
	ShopID           string       `json:"shopID"`
	AccessToken      string       `json:"accessToken"`
	Country          string       `json:"country"`
	RefreshToken     string       `json:"refreshToken"`
	AccountPlatform  string       `json:"accountPlatform"`
	Account          string       `json:"account"`
	UserIDVn         string       `json:"userIDVn"`
	SellerIDVn       string       `json:"sellerIDVn"`
	ShortCodeVn      string       `json:"shortCodeVn"`
	RefreshExpiresIn sql.NullTime `json:"refreshExpiresIn"`
	AccessExpiresIn  sql.NullTime `json:"accessExpiresIn"`
	Created          sql.NullTime `json:"created"`
}

type ProductsAttributeLazada struct {
	ID                   string         `json:"id"`
	Name                 sql.NullString `json:"name"`
	ShortDescription     sql.NullString `json:"shortDescription"`
	Description          sql.NullString `json:"description"`
	Brand                sql.NullString `json:"brand"`
	Model                sql.NullString `json:"model"`
	HeadphoneFeatures    sql.NullString `json:"headphoneFeatures"`
	Bluetooth            sql.NullString `json:"bluetooth"`
	WarrantyType         sql.NullString `json:"warrantyType"`
	Warranty             sql.NullString `json:"warranty"`
	Hazmat               sql.NullString `json:"hazmat"`
	ExpireDate           sql.NullString `json:"expireDate"`
	BrandClassification  sql.NullString `json:"brandClassification"`
	IngredientPreference sql.NullString `json:"ingredientPreference"`
	LotNumber            sql.NullString `json:"lotNumber"`
	UnitsHb              sql.NullString `json:"unitsHb"`
	FmltSkincare         sql.NullString `json:"fmltSkincare"`
	Quantitative         sql.NullString `json:"quantitative"`
	SkincareByAge        sql.NullString `json:"skincareByAge"`
	SkinBenefit          sql.NullString `json:"skinBenefit"`
	SkinType             sql.NullString `json:"skinType"`
	UserManual           sql.NullString `json:"userManual"`
	CountryOriginHb      sql.NullString `json:"countryOriginHb"`
	ColorFamily          sql.NullString `json:"colorFamily"`
	FragranceFamily      sql.NullString `json:"fragranceFamily"`
	Source               sql.NullString `json:"source"`
	ProductID            string         `json:"productID"`
}

type ProductsLazada struct {
	ID                    string         `json:"id"`
	LazadaID              int64          `json:"lazadaID"`
	LazadaPrimaryCategory int64          `json:"lazadaPrimaryCategory"`
	Created               time.Time      `json:"created"`
	Updated               time.Time      `json:"updated"`
	Status                sql.NullString `json:"status"`
	SubStatus             sql.NullString `json:"subStatus"`
	ShopID                string         `json:"shopID"`
}

type ProductsSkuLazada struct {
	ID              string         `json:"id"`
	Status          sql.NullString `json:"status"`
	Quantity        sql.NullInt32  `json:"quantity"`
	SellerSku       sql.NullString `json:"sellerSku"`
	ShopSku         sql.NullString `json:"shopSku"`
	Url             sql.NullString `json:"url"`
	ColorFamily     sql.NullString `json:"colorFamily"`
	Price           sql.NullInt32  `json:"price"`
	Available       sql.NullInt32  `json:"available"`
	SkuID           sql.NullInt64  `json:"skuID"`
	PackageContent  sql.NullString `json:"packageContent"`
	PackageWidth    sql.NullString `json:"packageWidth"`
	PackageWeight   sql.NullString `json:"packageWeight"`
	PackageLength   sql.NullString `json:"packageLength"`
	PackageHeight   sql.NullString `json:"packageHeight"`
	SpecialPrice    sql.NullString `json:"specialPrice"`
	SpecialToTime   sql.NullTime   `json:"specialToTime"`
	SpecialFromTime sql.NullTime   `json:"specialFromTime"`
	SpecialFromDate sql.NullTime   `json:"specialFromDate"`
	SpecialToDate   sql.NullTime   `json:"specialToDate"`
	ProductID       string         `json:"productID"`
}

type Shop struct {
	ID         string       `json:"id"`
	UserID     string       `json:"userID"`
	ShopName   string       `json:"shopName"`
	Created    sql.NullTime `json:"created"`
	LastUpdate sql.NullTime `json:"lastUpdate"`
}

type ShopsPlatform struct {
	ID           string       `json:"id"`
	ShopID       string       `json:"shopID"`
	PlatformName string       `json:"platformName"`
	Created      sql.NullTime `json:"created"`
}

type User struct {
	ID       string       `json:"id"`
	Password string       `json:"password"`
	Email    string       `json:"email"`
	Created  sql.NullTime `json:"created"`
	Fullname string       `json:"fullname"`
}
