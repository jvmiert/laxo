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
    el := models.Element{Type: nodeType, Children: []models.Text{}}

    if align != "" {
      el.Align = align
    }

		schema = append(schema, el)
		index++
	}

	return schema, index
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
        }
				prevNode = "div"
			case "span":
				//@HACK: Because Lazada is retarded.
				if prevNode == "div" {
					schema, index = addNode(schema, "paragraph", index, "")
					prevNode = "span"
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
							Children: []models.Text{},
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
							schema[index].Children = []models.Text{}
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
							Children: []models.Text{},
							Width:    asset.Width.Int64,
							Height:   asset.Height.Int64,
						}

						schema = append(schema, el)
						index++
					}
				}
			}
		case html.TextToken:
			trimmed := strings.TrimSpace(token.String())
			if trimmed != "" {
				//s.server.Logger.Debugw("TextToken", "token", token, "trimmed", trimmed, "depth", depth, "index", index)
				if depth > 0 && index != -1 {
					schema[index].Children = append(schema[index].Children, models.Text{Text: trimmed, Bold: bold, Underline: underline, Italic: italic})
				}
			}
		}
	}

	if len(schema) == 0 {
		text := s.GetSantizedString(h)
		if text == "" {
			return "", ErrEmptyText
		}
		schema = append(schema, models.Element{Type: "paragraph", Children: []models.Text{{Text: text}}})
	}

	b, _ := json.Marshal(schema)

	return string(b), nil
}
