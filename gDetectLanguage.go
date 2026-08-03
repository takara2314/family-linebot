package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate/apiv3/translatepb"
	"golang.org/x/text/language"
)

func gDetectLanguage(text string) language.Tag {
	ctx := context.Background()
	response, err := translateClient.DetectLanguage(ctx, &translatepb.DetectLanguageRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/global", projectID),
		Source:   &translatepb.DetectLanguageRequest_Content{Content: text},
		MimeType: "text/plain",
	})
	if err != nil {
		log.Printf("DetectLanguage: %v", err)
		panic(err)
	}
	if len(response.Languages) == 0 {
		panic("DetectLanguage returned empty response")
	}
	return language.Make(response.Languages[0].LanguageCode)
}
