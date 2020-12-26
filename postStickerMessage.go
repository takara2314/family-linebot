package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func postStickerMessage(event *linebot.Event) {
	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage("こんにちは！"),
	).Do()

	if err != nil {
		log.Println(err)
		panic(err)
	}
}
