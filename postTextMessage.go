package main

import "github.com/line/line-bot-sdk-go/linebot"

func postTextMessage(event *http.Event) {
	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage("こんにちは！")
	).Do()

	if _, err != nil {
		panic(err)
	}
}