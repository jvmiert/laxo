package shop

import (
	"strings"

	"golang.org/x/net/html"
)

func addNode(schema []Element, nodeType string, index int) ([]Element, int) {
	add := false

	if len(schema) > 0 {
		if len(schema[len(schema)-1].Children) > 0 {
			add = true
		} else {
			schema[len(schema)-1].Type = nodeType
		}
	} else {
		add = true
	}

	if add {
		schema = append(schema, Element{Type: nodeType, Children: []Text{}})
		index++
	}

	return schema, index
}

func (s *Service) HTMLToSlate(h string) ([]Element, error) {
	tkn := html.NewTokenizer(strings.NewReader(h))

	schema := []Element{}
	index := -1
	depth := 0
	prevNode := ""

	bold := false
	italic := false
	underline := false

out:
	for {
		tt := tkn.Next()
		token := tkn.Token()
		switch tt {
		case html.ErrorToken:
			if len(schema[len(schema)-1].Children) == 0 {
				schema = schema[:len(schema)-1]
			}
			s.server.Logger.Debug("ErrorToken")
			break out
		case html.StartTagToken:
			s.server.Logger.Debugw("StartTagToken", "token", token, "type", token.Data)
			switch token.Data {
			case "div":
				prevNode = "div"
			case "span":
				//@HACK: Because Lazada is retarded.
				if prevNode == "div" {
					schema, index = addNode(schema, "paragraph", index)
					prevNode = "span"
				}
				depth++
			case "p":
				if depth == 0 {
					schema, index = addNode(schema, "paragraph", index)
					prevNode = "p"
				}
				depth++
			case "h1":
				if depth == 0 {
					schema, index = addNode(schema, "heading-one", index)
					prevNode = "h1"
				}
				depth++
			case "h2":
				if depth == 0 {
					schema, index = addNode(schema, "heading-two", index)
					prevNode = "h2"
				}
				depth++
			case "h3":
				if depth == 0 {
					schema, index = addNode(schema, "heading-three", index)
					prevNode = "h3"
				}
				depth++
			case "strong":
				bold = true
			case "em":
				italic = true
			case "u":
				underline = true
			}
		case html.EndTagToken:
			s.server.Logger.Debugw("EndTagToken", "token", token, "type", token.Data)
			switch token.Data {
			case "span":
				depth--
			case "p":
				depth--
			case "h1":
				depth--
			case "h2":
				depth--
			case "h3":
				depth--
			case "strong":
				bold = false
			case "em":
				italic = false
			case "u":
				underline = false
			}
		case html.SelfClosingTagToken:
			s.server.Logger.Debugw("SelfClosingTagToken", "token", token)
		case html.TextToken:
			trimmed := strings.TrimSpace(token.String())
			if trimmed != "" {
				s.server.Logger.Debugw("TextToken", "token", token, "trimmed", trimmed, "depth", depth, "index", index)
				if depth > 0 {
					schema[index].Children = append(schema[index].Children, Text{Text: trimmed, Bold: bold, Underline: underline, Italic: italic})
				}
			}
		}
	}

	return schema, nil
}
