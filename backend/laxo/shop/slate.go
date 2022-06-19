package shop

type Text struct {
	Text      string `json:"text"`
	Bold      bool   `json:"bold,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	Italic    bool   `json:"italic,omitempty"`
}

type Element struct {
	Type     string `json:"type,omitempty"`
	Src      string `json:"src,omitempty"`
	Width    int64  `json:"width,omitempty"`
	Height   int64  `json:"height,omitempty"`
	Children []Text `json:"children"`
}
