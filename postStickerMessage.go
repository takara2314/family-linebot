package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"google.golang.org/genai"
)

func postStickerMessage(event *linebot.Event, stickerID string) {
	ctx := context.Background()
	url := "https://stickershop.line-scdn.net/stickershop/v1/sticker/" + stickerID + "/android/sticker.png"
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	image, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	result, err := geminiClient.Models.GenerateContent(ctx, appConfig.GeminiStickerModel, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			genai.NewPartFromText("この画像はLINEスタンプです。画像全体の構造、文字の配置、吹き出し、キャラクターの表情を考慮して、実際に見える文字を読み取ってください。画像にない言葉は創作せず、読み順と表記だけを通常の発話として整えてください。"),
			genai.NewPartFromBytes(image, "image/png"),
		},
	}}, &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"has_visible_text":   {Type: genai.TypeBoolean},
				"normalized_message": {Type: genai.TypeString},
				"uncertain":          {Type: genai.TypeBoolean},
			},
			Required: []string{"has_visible_text", "normalized_message", "uncertain"},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	var detected struct {
		HasVisibleText    bool   `json:"has_visible_text"`
		NormalizedMessage string `json:"normalized_message"`
		Uncertain         bool   `json:"uncertain"`
	}
	if err := json.Unmarshal([]byte(result.Text()), &detected); err != nil {
		log.Println(err)
		return
	}
	if !detected.HasVisibleText || detected.Uncertain || detected.NormalizedMessage == "" {
		return
	}

	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(convertJpTh(detected.NormalizedMessage))).Do()
	if err != nil {
		log.Println(fmt.Errorf("reply sticker: %w", err))
		panic(err)
	}
}
