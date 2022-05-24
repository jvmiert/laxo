package shop

import (
	"encoding/json"
	"sort"
)

type OAuthVerifyRequest struct {
	Platform string `json:"platform"`
	Code     string `json:"code"`
	State    string `json:"state"`
}

type ReturnRedirect struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

type availablePlatforms map[string]struct{}

func (a availablePlatforms) Has(v string) bool {
	_, ok := a[v]
	return ok
}

type OAuthRedirectRequest struct {
	ShopID          string            `json:"shopID"`
	ReturnRedirects []*ReturnRedirect `json:"platforms"`
  Connected       []string          `json:"connectedPlatforms"`
}

func (s *OAuthRedirectRequest) JSON() ([]byte, error) {
	sort.Slice(s.ReturnRedirects, func(i, j int) bool { return s.ReturnRedirects[i].Platform < s.ReturnRedirects[j].Platform })
	bytes, err := json.Marshal(s)

	if err != nil {
		return bytes, err
	}

	return bytes, nil
}
