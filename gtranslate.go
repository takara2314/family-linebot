package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func gtranslate(text string) string {
	ctx := context.Background()

	lang := language.English

	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
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
