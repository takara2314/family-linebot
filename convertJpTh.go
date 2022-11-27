package main

import (
	"fmt"

	"golang.org/x/text/language"
)

func convertJpTh(text string) string {
	var detectLang language.Tag = gDetectLanguage(text)

	switch detectLang {
	case language.Japanese:
		return gTranslate(text, language.Thai)

	case language.Thai:
		return gTranslate(text, language.Japanese)

	default:
		return fmt.Sprintf(
			"%s\n(%s)",
			gTranslate(text, language.Japanese),
			gTranslate(text, language.Thai),
		)
	}
}
