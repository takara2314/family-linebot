package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot *linebot.Client
	err error
)

func main() {
	bot, err = linebot.New(
		os.Getenv("LINEBOT_CHANNEL_SECRET"),
		os.Getenv("LINEBOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/", rootGET)
	router.POST("/callback", callbackPOST)

	router.Run(":" + os.Getenv("PORT"))
}

func rootGET(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func callbackPOST(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch event.Message.(type) {
			case *linebot.TextMessage:
				postTextMessage(event)

			case *linebot.StickerMessage:
				postStickerMessage(event)
			}
		}
	}
}
