package main

import (
	"context"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func gDetectLanguage(text string) language.Tag {
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Printf("translate.NewClient: %v", err)
		panic(err)
	}
	defer client.Close()

	lang, err := client.DetectLanguage(ctx, []string{text})
	if err != nil {
		log.Printf("DetectLanguage: %v", err)
		panic(err)
	}
	if len(lang) == 0 || len(lang[0]) == 0 {
		log.Printf("DetectLanguage return value empty")
		panic("DetectLanguage return value empty")
	}
	return lang[0][0].Language
}
