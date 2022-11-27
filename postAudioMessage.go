package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/line/line-bot-sdk-go/linebot"
)

func postAudioMessage(event *linebot.Event, messageID string) {
	ctx := context.Background()

	// get audio content
	content, err := bot.GetMessageContent(messageID).Do()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer content.Content.Close()

	// convert audio content to bytes
	buf := new(bytes.Buffer)
	io.Copy(buf, content.Content)
	ret := buf.Bytes()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer client.Close()

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_WEBM_OPUS,
			SampleRateHertz: 48000,
			LanguageCode:    "ja-JP",
			Model:           "default",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: ret},
		},
	})
	if err != nil {
		fmt.Printf("failed to recognize: %v", err)
		panic(err)
	}

	replyMessage := ""

	for _, result := range resp.Results {
		replyMessage += result.Alternatives[0].Transcript
	}

	if replyMessage == "" {
		replyMessage = "[エラー] 音声を認識できませんでした。 (เราจำเสียงของคุณไม่ได้)"
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
