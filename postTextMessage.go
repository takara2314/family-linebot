package main

import (
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/text/language"
)

func postTextMessage(event *linebot.Event, message string) {
	var replyMessage string
	var detectLang language.Tag = gDetectLanguage(message)

	switch detectLang {
	case language.Japanese:
		replyMessage = gTranslate(message, language.Thai)

	case language.Thai:
		replyMessage = gTranslate(message, language.Japanese)

	default:
		replyMessage = fmt.Sprintf(
			"%s\n(%s)",
			gTranslate(message, language.Japanese),
			gTranslate(message, language.Thai),
		)
	}

	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(replyMessage),
	).Do()

	if err != nil {
		log.Println(err)
		panic(err)
	}
}
