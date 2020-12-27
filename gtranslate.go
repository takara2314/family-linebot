package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func gTranslate(text string, targetLang language.Tag) string {
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, targetLang, nil)
	if err != nil {
		fmt.Printf("Translate: %v", err)
		panic(err)
	}
	if len(resp) == 0 {
		fmt.Printf("Translate returned empty response to text: %s", text)
		panic(err)
	}
	return resp[0].Text
}
