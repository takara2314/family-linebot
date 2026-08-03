package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/translate/apiv3/translatepb"
	"golang.org/x/text/language"
)

func gTranslate(text string, targetLang language.Tag) string {
	ctx := context.Background()
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	contents := make([]string, 0, len(lines))
	indexes := make([]int, 0, len(lines))
	prefixes := make([]string, len(lines))
	suffixes := make([]string, len(lines))

	for i, line := range lines {
		withoutPrefix := strings.TrimLeft(line, " \t")
		prefixes[i] = line[:len(line)-len(withoutPrefix)]
		content := strings.TrimRight(withoutPrefix, " \t")
		suffixes[i] = withoutPrefix[len(content):]
		if content == "" {
			continue
		}
		contents = append(contents, content)
		indexes = append(indexes, i)
	}
	if len(contents) == 0 {
		return text
	}

	response, err := translateClient.TranslateText(ctx, &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		Contents:           contents,
		MimeType:           "text/plain",
		TargetLanguageCode: targetLang.String(),
	})
	if err != nil {
		log.Printf("Translate: %v", err)
		panic(err)
	}
	if len(response.Translations) != len(contents) {
		panic("Translate returned unexpected number of results")
	}
	for i, translated := range response.Translations {
		index := indexes[i]
		lines[index] = prefixes[index] + translated.TranslatedText + suffixes[index]
	}
	return strings.Join(lines, "\n")
}
