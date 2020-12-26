package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func postTextMessage(event *linebot.Event, message string) {
	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(gtranslate(message)),
	).Do()

	if err != nil {
		log.Println(err)
		panic(err)
	}
}
