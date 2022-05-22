package laxo

import (
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	_ "laxo.vn/laxo/laxo/translations"
)

var serverLangs = []language.Tag{
    language.English,
    language.Vietnamese,
}

var matcher = language.NewMatcher(serverLangs)

func GetLocalePrinter(r *http.Request) *message.Printer {
  localeHeader := r.Header.Get("locale")

  localeTag, _, _ := matcher.Match(language.Make(localeHeader))

  return message.NewPrinter(localeTag)
}

