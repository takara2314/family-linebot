package main

import "github.com/line/line-bot-sdk-go/linebot"

func postTextMessage(event *linebot.Event) {
	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage("こんにちは！"),
	).Do()

	if err != nil {
		panic(err)
	}
}
