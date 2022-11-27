package main

import (
	"context"
	"log"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/line/line-bot-sdk-go/linebot"
)

func postStickerMessage(event *linebot.Event, stickerID string) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	texts, err := client.DetectTexts(
		ctx,
		vision.NewImageFromURI("https://stickershop.line-scdn.net/stickershop/v1/sticker/"+stickerID+"/android/sticker.png"),
		nil,
		10,
	)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if len(texts) > 0 {
		detectedText := texts[0].Description

		detectedText = strings.Replace(detectedText, " ", "", -1)
		detectedText = strings.Replace(detectedText, "　", "", -1)
		detectedText = strings.Replace(detectedText, "\n", "", -1)
		detectedText = strings.Replace(detectedText, ".", "", -1)

		_, err = bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(convertJpTh(detectedText)),
		).Do()

		if err != nil {
			log.Println(err)
			panic(err)
		}
	}
}
