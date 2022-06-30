package models

type Paginate struct {
	Total  int64 `json:"total"`
	Pages  int64 `json:"pages"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}
