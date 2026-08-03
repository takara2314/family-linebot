package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"cloud.google.com/go/speech/apiv2/speechpb"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

func postAudioMessage(event *linebot.Event, messageID string) {
	ctx := context.Background()
	content, err := bot.GetMessageContent(messageID).Do()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer content.Content.Close()
	audio, err := io.ReadAll(content.Content)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	response, err := speechClient.Recognize(ctx, &speechpb.RecognizeRequest{
		Recognizer: fmt.Sprintf("projects/%s/locations/global/recognizers/_", projectID),
		Config: &speechpb.RecognitionConfig{
			DecodingConfig: &speechpb.RecognitionConfig_AutoDecodingConfig{AutoDecodingConfig: &speechpb.AutoDetectDecodingConfig{}},
			Model:          "chirp_3",
			LanguageCodes:  []string{"ja-JP"},
		},
		AudioSource: &speechpb.RecognizeRequest_Content{Content: audio},
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	var replyMessage strings.Builder
	for _, result := range response.Results {
		if len(result.Alternatives) > 0 {
			replyMessage.WriteString(result.Alternatives[0].Transcript)
		}
	}
	if replyMessage.Len() == 0 {
		replyMessage.WriteString("[エラー] 音声を日本語として認識できませんでした。")
	}
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage.String())).Do()
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
