package main

import (
	"context"
	"log"
	"net/http"

	speech "cloud.google.com/go/speech/apiv2"
	translate "cloud.google.com/go/translate/apiv3"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"google.golang.org/api/option"
	"google.golang.org/genai"
)

var (
	bot             *linebot.Client
	translateClient *translate.TranslationClient
	speechClient    *speech.Client
	geminiClient    *genai.Client
	projectID       string
	appConfig       config
	err             error
)

func main() {
	ctx := context.Background()
	appConfig, err = loadConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	projectID = appConfig.ProjectID

	bot, err = linebot.New(appConfig.LineChannelSecret, appConfig.LineChannelToken)
	if err != nil {
		log.Fatal(err)
	}
	translateClient, err = translate.NewTranslationClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer translateClient.Close()
	speechClient, err = speech.NewClient(ctx, option.WithEndpoint(appConfig.SpeechLocation+"-speech.googleapis.com:443"))
	if err != nil {
		log.Fatal(err)
	}
	defer speechClient.Close()
	geminiClient, err = genai.NewClient(ctx, &genai.ClientConfig{
		Backend:  genai.BackendEnterprise,
		Project:  projectID,
		Location: appConfig.GeminiLocation,
	})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/", rootGET)
	router.POST("/callback", callbackPOST)
	if err := router.Run(":" + appConfig.Port); err != nil {
		log.Fatal(err)
	}
}

func rootGET(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func callbackPOST(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(http.StatusBadRequest)
		} else {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				postTextMessage(event, message.Text)
			case *linebot.StickerMessage:
				postStickerMessage(event, message.StickerID)
			case *linebot.AudioMessage:
				postAudioMessage(event, message.ID)
			}
		}
	}
}
