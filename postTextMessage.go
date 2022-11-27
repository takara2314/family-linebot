package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func postTextMessage(event *linebot.Event, message string) {
	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(convertJpTh(message)),
	).Do()

	if err != nil {
		log.Println(err)
		panic(err)
	}
}
