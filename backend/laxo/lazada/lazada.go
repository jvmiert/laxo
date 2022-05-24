package lazada

import "gopkg.in/guregu/null.v4"

type PlatformLazadaReturn struct {
	ID               string    `json:"id"`
	ShopID           string    `json:"shopID"`
	Country          string    `json:"country"`
	AccountPlatform  string    `json:"accountPlatform"`
	Account          string    `json:"account"`
	UserIDVn         string    `json:"userIDVn"`
	SellerIDVn       string    `json:"sellerIDVn"`
	ShortCodeVn      string    `json:"shortCodeVn"`
	RefreshExpiresIn null.Time `json:"refreshExpiresIn"`
	AccessExpiresIn  null.Time `json:"accessExpiresIn"`
	Created          null.Time `json:"created"`
}
