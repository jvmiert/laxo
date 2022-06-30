package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/models"
	"laxo.vn/laxo/laxo/shop"
	"laxo.vn/laxo/laxo/store"
)

func TestSlateToHTML(t *testing.T) {
	os.Chdir("./..")

	if err := godotenv.Load(".env"); err != nil {
		t.Errorf("Failed to load .env file: %v", err)
	}

	logger := laxo.NewLogger()
	defer logger.Zap.Sync()

	config, err := laxo.InitConfig()
	if err != nil {
		t.Errorf("Could not init config: %v", err)
	}

	server, err := laxo.NewServer(logger, config)
	if err != nil {
		t.Errorf("Failed to get server struct: %v", err)
	}

	dbURI := os.Getenv("POSTGRESQL_TEST_URL")

	if err = server.InitDatabase(dbURI); err != nil {
		logger.Errorw("Failed to init Database",
			"error", err,
			"uri", dbURI,
		)
		return
	}

	assetsBasePath := os.Getenv("ASSETS_BASE_PATH")
	store, err := store.NewStore(dbURI, logger, assetsBasePath)
	if err != nil {
		logger.Error("Failed to create new store", "error", err)
		return
	}

	shopService := shop.NewService(store, logger, server)

	tests := []struct {
		id       string
		html     string
		slate    []models.Element
		laxoHTML string
	}{
		{
			"Lazada simple text",
			`<p style="margin:0"><span>This is left aligned</span></p>`,
			[]models.Element{models.Element{Type: "paragraph", Children: []models.Element{models.Element{Text: "This is left aligned"}}}},
			`<p>This is left aligned</p>`,
		},
		{
			"Lazada center align",
			`<p style="text-align:center;display:inline-block;width:100%;margin:0"><span>This is center aligned</span></p>`,
			[]models.Element{models.Element{Type: "paragraph", Align: "center", Children: []models.Element{models.Element{Text: "This is center aligned"}}}},
			`<p style="text-align: center;">This is center aligned</p>`,
		},
		{
			"Lazada right align",
			`<p style="text-align:right;display:inline-block;width:100%;margin:0"><span>This is right aligned</span></p>`,
			[]models.Element{models.Element{Type: "paragraph", Align: "right", Children: []models.Element{models.Element{Text: "This is right aligned"}}}},
			`<p style="text-align: right;">This is right aligned</p>`,
		},
		{
			"Lazada ul list",
			`<ul style="list-style:disc;margin-left:10px">
				 <li><span>This is bullet point one</span></li>
				 <li><span style="text-align:right;display:inline-block;width:100%">This is bullet point two, right aligned and <strong style="font-weight:bold">bold</strong></span></li>
				 <li><span>This is bullet point three</span></li>
			</ul>`,
			[]models.Element{models.Element{Type: "bulleted-list", Children: []models.Element{
				models.Element{Type: "list-item", Children: []models.Element{
					models.Element{Text: "This is bullet point one"},
				}},
				models.Element{Type: "list-item", Align: "right", Children: []models.Element{
					models.Element{Text: "This is bullet point two, right aligned and "},
					models.Element{Text: "bold", Bold: true},
				}},
				models.Element{Type: "list-item", Children: []models.Element{
					models.Element{Text: "This is bullet point three"},
				}},
			}}},
			`<ul><li>This is bullet point one</li><li style="text-align: right;">This is bullet point two, right aligned and <strong>bold</strong></li><li>This is bullet point three</li></ul>`,
		},
	}

	for _, v := range tests {
		t.Run(v.id, func(t *testing.T) {
			schema, err := shopService.HTMLToSlate(v.html, "01G1FZCVYH9J47DB2HZENSBC6E")
			if err != nil {
				t.Errorf("HTMLToSlate returned error: %v", err)
			}

			b, err := json.Marshal(v.slate)
			if err != nil {
				t.Errorf("Marshal returned error: %v", err)
			}

			assert.Equal(t, schema, string(b), "generated schema does not match")

			encodedSchema := &[]models.Element{}
			err = json.Unmarshal([]byte(schema), encodedSchema)
			if err != nil {
				t.Errorf("Unmarshal returned error: %v", err)
			}

			html, err := shopService.SlateToHTML(*encodedSchema)
			if err != nil {
				t.Errorf("Marshal returned error: %v", err)
			}

			assert.Equal(t, v.laxoHTML, html, "generated html does not match")
		})
	}
}
