package shop

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/tdewolff/parse/css"
	"golang.org/x/net/html"
	"laxo.vn/laxo/laxo/models"
)

var ErrEmptyText = errors.New("html parse resulted in empty text")

// We're just hacking it up right now to do a half acceptable job in parsing
// Lazada's garbage.
func parseCSS(token *html.Token) (alignValue string, isInlineBlock bool) {
	var style string

	for _, a := range token.Attr {
		if a.Key == "style" {
			style = a.Val
		}
	}

	if style != "" {
		p := css.NewParser(bytes.NewBufferString(style), true)
	out:
		for {
			gt, _, data := p.Next()
			switch gt {
			case css.DeclarationGrammar:
				values := p.Values()
				declarationValue := string(data)
				if declarationValue == "text-align" {
					if len(values) == 1 {
						alignValue = string(values[0].Data)
					}
				}
				if declarationValue == "display" {
					if len(values) == 1 {
						displayValue := string(values[0].Data)

						if displayValue == "inline-block" {
							isInlineBlock = true
						}
					}
				}
			case css.ErrorGrammar:
				break out
			}
		}
	}
	return alignValue, isInlineBlock
}

func addNode(schema []models.Element, nodeType string, index int, align string) ([]models.Element, int) {
	add := false

	if len(schema) > 0 {
		if len(schema[len(schema)-1].Children) > 0 || schema[len(schema)-1].Type == "image" {
			add = true
		} else {
			schema[len(schema)-1].Type = nodeType
			if align != "" {
				schema[len(schema)-1].Align = align
			}
		}
	} else {
		add = true
	}

	if add {
		el := models.Element{Type: nodeType, Children: []models.Element{}}

		if align != "" {
			el.Align = align
		}

		schema = append(schema, el)
		index++
	}

	return schema, index
}

func getAlignStyle(align string) string {
	var alignStyle string

	switch align {
	case "left":
		alignStyle = ` style="text-align: left;"`
	case "right":
		alignStyle = ` style="text-align: right;"`
	case "center":
		alignStyle = ` style="text-align: center;"`
	case "justify":
		alignStyle = ` style="text-align: justify;white-space: pre-line;"`
	}

	return alignStyle
}

func slateElementToHTML(element models.Element) string {
	var innerHtml string
	var html string

	for _, innerNode := range element.Children {
		// opening styling tags
		if innerNode.Bold {
			innerHtml = innerHtml + "<strong>"
		}

		if innerNode.Underline {
			innerHtml = innerHtml + "<u>"
		}

		if innerNode.Italic {
			innerHtml = innerHtml + "<i>"
		}

		// actual text
		innerHtml = innerHtml + innerNode.Text

		// closing styling tags
		if innerNode.Bold {
			innerHtml = innerHtml + "</strong>"
		}

		if innerNode.Underline {
			innerHtml = innerHtml + "</u>"
		}

		if innerNode.Italic {
			innerHtml = innerHtml + "</i>"
		}
	}

	alignStyle := getAlignStyle(element.Align)

	//@TODO: add support for images
	switch element.Type {
	case "paragraph":
		html = "<p" + alignStyle + ">" + innerHtml + "</p>"
	case "heading-one":
		html = "<h1" + alignStyle + ">" + innerHtml + "</h1>"
	case "heading-two":
		html = "<h2" + alignStyle + ">" + innerHtml + "</h2>"
	case "heading-three":
		html = "<h3" + alignStyle + ">" + innerHtml + "</h3>"
	case "list-item":
		html = "<li" + alignStyle + ">" + innerHtml + "</li>"
	default:
		html = innerHtml
	}

	fmt.Println("slateElementToHTML", "align", element.Align)
	//fmt.Println("slateElementToHTML", "html", html, "element", element)
	return html
}

func (s *Service) SlateToHTML(slateSchema []models.Element) (string, error) {
	var html string
	for _, node := range slateSchema {
		if node.Type == "numbered-list" || node.Type == "bulleted-list" {
			if node.Type == "bulleted-list" {
				html = html + "<ul>"
			}
			if node.Type == "numbered-list" {
				html = html + "<ol>"
			}
			// These are the <li> elements
			for _, listNode := range node.Children {
				html = html + slateElementToHTML(listNode)
			}
			if node.Type == "bulleted-list" {
				html = html + "</ul>"
			}
			if node.Type == "numbered-list" {
				html = html + "</ol>"
			}
		} else {
			html = html + slateElementToHTML(node)
		}
	}
	return html, nil
}

func (s *Service) HTMLToSlate(h string, shopID string) (string, error) {
	tkn := html.NewTokenizer(strings.NewReader(h))

	schema := []models.Element{}
	index := -1
	depth := 0
	prevNode := ""

	bold := false
	italic := false
	underline := false

	addToList := false

out:
	for {
		tt := tkn.Next()
		token := tkn.Token()
		switch tt {
		case html.ErrorToken:
			if index == -1 {
				break out
			}
			if len(schema[len(schema)-1].Children) == 0 && schema[len(schema)-1].Type != "image" {
				schema = schema[:len(schema)-1]
			}
			//s.server.Logger.Debug("ErrorToken")
			break out
		case html.StartTagToken:
			//s.server.Logger.Debugw("StartTagToken", "token", token, "type", token.Data, "depth", depth, "index", index)
			switch token.Data {
			case "div":
				align, isInlineBlock := parseCSS(&token)
				//@HACK: Lazada does dumb shit (just wrapping text with a div instead of adding a span like they normally do) with text aligning
				if isInlineBlock {
					schema, index = addNode(schema, "paragraph", index, align)
					depth++
					prevNode = "inlineBlockDiv"
				} else {
					prevNode = "div"
				}
			case "span":
				//@HACK: Because Lazada is retarded.
				if prevNode == "div" {
					schema, index = addNode(schema, "paragraph", index, "")
					prevNode = "span"
				}
				//@HACK: More lazada retardness (using a span inside a <li> element to align)
				if addToList {
					align, _ := parseCSS(&token)
					if align != "" {
						listItemIndex := len(schema[index].Children) - 1
						schema[index].Children[listItemIndex].Align = align
					}
				}
				depth++
			case "p":
				if depth == 0 {
					align, _ := parseCSS(&token)
					schema, index = addNode(schema, "paragraph", index, align)
					prevNode = "p"
				}
				depth++
			case "h1":
				if depth == 0 {
					schema, index = addNode(schema, "heading-one", index, "")
					prevNode = "h1"
				}
				depth++
			case "h2":
				if depth == 0 {
					schema, index = addNode(schema, "heading-two", index, "")
					prevNode = "h2"
				}
				depth++
			case "h3":
				if depth == 0 {
					schema, index = addNode(schema, "heading-three", index, "")
					prevNode = "h3"
				}
				depth++
			case "ul":
				if depth == 0 {
					schema, index = addNode(schema, "bulleted-list", index, "")
					prevNode = "ul"
				}
				depth++
			case "ol":
				if depth == 0 {
					schema, index = addNode(schema, "numbered-list", index, "")
					prevNode = "ol"
				}
				depth++
			case "li":
				if schema[index].Type == "numbered-list" || schema[index].Type == "bulleted-list" {
					align, _ := parseCSS(&token)
					schema[index].Children = append(schema[index].Children, models.Element{Type: "list-item", Align: align, Children: []models.Element{}})
					addToList = true
				}
				//el := models.Element{Type: "list-item", Children: []models.Element{}}
				//schema, index = addNode(schema, "list-item", index, "")
				//prevNode = "li"
			case "img":
				if depth == 0 {
					findImages := regexp.MustCompile(`src=["'](.*?)["']`)
					matches := findImages.FindStringSubmatch(token.String())

					if len(matches) > 1 {
						//s.server.Logger.Debugw("ADDING an image", "data", token, "src", matches[1])

						originalName := path.Base(matches[1])

						asset, err := s.store.GetAssetByOriginalName(originalName, shopID)
						if err != pgx.ErrNoRows && err != nil {
							return "", fmt.Errorf("GetAssetByOriginalName: %w", err)
						}

						if errors.Is(err, pgx.ErrNoRows) {
							continue
						}

						el := models.Element{
							Type:     "image",
							Src:      asset.ID + asset.Extension.String,
							Children: []models.Element{},
							Width:    asset.Width.Int64,
							Height:   asset.Height.Int64,
						}

						schema = append(schema, el)
						index++
					}
				}

				if prevNode == "p" {
					if len(schema[index].Children) == 0 {
						findImages := regexp.MustCompile(`src=["'](.*?)["']`)
						matches := findImages.FindStringSubmatch(token.String())

						if len(matches) > 1 {
							//s.server.Logger.Debugw("ADDING an image", "data", token, "src", matches[1])

							originalName := path.Base(matches[1])

							asset, err := s.store.GetAssetByOriginalName(originalName, shopID)
							if err != pgx.ErrNoRows && err != nil {
								return "", fmt.Errorf("GetAssetByOriginalName: %w", err)
							}

							if errors.Is(err, pgx.ErrNoRows) {
								continue
							}

							schema[index].Type = "image"
							schema[index].Src = asset.ID + asset.Extension.String
							schema[index].Children = []models.Element{}
							schema[index].Width = asset.Width.Int64
							schema[index].Height = asset.Height.Int64
						}
					}
				}
			case "strong":
				bold = true
			case "em":
				italic = true
			case "u":
				underline = true
			}
		case html.EndTagToken:
			//s.server.Logger.Debugw("EndTagToken", "token", token, "type", token.Data)
			switch token.Data {
			case "div":
				if prevNode == "inlineBlockDiv" {
					depth--
					prevNode = ""
				}
			case "span":
				depth--
			case "p":
				depth--
			case "h1", "h2", "h3":
				depth--
			case "ol", "ul":
				depth--
			case "li":
				addToList = false
			case "strong":
				bold = false
			case "em":
				italic = false
			case "u":
				underline = false
			}
		case html.SelfClosingTagToken:
			//s.server.Logger.Debugw("SelfClosingTagToken", "token", token, "type", token.Data, "depth", depth, "index", index)
			switch token.Data {
			case "img":
				if depth == 0 {
					findImages := regexp.MustCompile(`src=["'](.*?)["']`)
					matches := findImages.FindStringSubmatch(token.String())

					if len(matches) > 1 {
						//s.server.Logger.Debugw("ADDING an image", "data", token, "src", matches[1])

						originalName := path.Base(matches[1])

						asset, err := s.store.GetAssetByOriginalName(originalName, shopID)
						if err != pgx.ErrNoRows && err != nil {
							return "", fmt.Errorf("GetAssetByOriginalName: %w", err)
						}

						if errors.Is(err, pgx.ErrNoRows) {
							continue
						}

						el := models.Element{
							Type:     "image",
							Src:      asset.ID + asset.Extension.String,
							Children: []models.Element{},
							Width:    asset.Width.Int64,
							Height:   asset.Height.Int64,
						}

						schema = append(schema, el)
						index++
					}
				}
			}
		case html.TextToken:
			tokenText := token.String()
			if tokenText != "" {
				//s.server.Logger.Debugw("TextToken", "token", token, "trimmed", trimmed, "depth", depth, "index", index)
				if addToList {
					listItemIndex := len(schema[index].Children) - 1
					schema[index].Children[listItemIndex].Children = append(schema[index].Children[listItemIndex].Children, models.Element{Text: tokenText, Bold: bold, Underline: underline, Italic: italic})
				} else if depth > 0 && index != -1 {
					schema[index].Children = append(schema[index].Children, models.Element{Text: tokenText, Bold: bold, Underline: underline, Italic: italic})
				}
			}
		}
	}

	if len(schema) == 0 {
		text := s.GetSantizedString(h)
		if text == "" {
			return "", ErrEmptyText
		}
		schema = append(schema, models.Element{Type: "paragraph", Children: []models.Element{{Text: text}}})
	}

	b, _ := json.Marshal(schema)

	return string(b), nil
}
