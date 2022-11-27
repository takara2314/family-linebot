package main

import (
	"context"
	"fmt"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/line/line-bot-sdk-go/linebot"
)

func postAudioMessage(event *linebot.Event, audioURL string) {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer client.Close()

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: audioURL},
		},
	})
	if err != nil {
		fmt.Printf("failed to recognize: %v", err)
		panic(err)
	}

	replyMessage := "<文字起こし> "

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			replyMessage += alt.Transcript
		}
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
