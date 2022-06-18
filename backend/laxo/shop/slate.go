package shop

type Text struct {
	Text      string `json:"text,omitempty"`
	Bold      bool   `json:"bold,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	Italic    bool   `json:"italic,omitempty"`
}

type Element struct {
	Type     string `json:"type,omitempty"`
	Children []Text `json:"children"`
}
