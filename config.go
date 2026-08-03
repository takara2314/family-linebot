package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type config struct {
	ProjectID          string
	Port               string
	LineChannelSecret  string
	LineChannelToken   string
	GeminiLocation     string
	GeminiStickerModel string
	SpeechLocation     string
}

func loadConfig(ctx context.Context) (config, error) {
	cfg := config{
		ProjectID:          os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Port:               envOrDefault("PORT", "8080"),
		LineChannelSecret:  os.Getenv("LINEBOT_CHANNEL_SECRET"),
		LineChannelToken:   os.Getenv("LINEBOT_CHANNEL_TOKEN"),
		GeminiLocation:     envOrDefault("GEMINI_LOCATION", "global"),
		GeminiStickerModel: envOrDefault("GEMINI_STICKER_MODEL", "gemini-3.5-flash-lite"),
		SpeechLocation:     envOrDefault("SPEECH_LOCATION", "us"),
	}
	if cfg.ProjectID == "" {
		return config{}, fmt.Errorf("GOOGLE_CLOUD_PROJECT is required")
	}
	if cfg.LineChannelSecret != "" && cfg.LineChannelToken != "" {
		return cfg, nil
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return config{}, fmt.Errorf("create Secret Manager client: %w", err)
	}
	defer client.Close()

	if cfg.LineChannelSecret == "" {
		cfg.LineChannelSecret, err = getSecret(ctx, client, cfg.ProjectID, envOrDefault("LINE_CHANNEL_SECRET_NAME", "linebot-channel-secret"))
		if err != nil {
			return config{}, err
		}
	}
	if cfg.LineChannelToken == "" {
		cfg.LineChannelToken, err = getSecret(ctx, client, cfg.ProjectID, envOrDefault("LINE_CHANNEL_TOKEN_NAME", "linebot-channel-token"))
		if err != nil {
			return config{}, err
		}
	}
	return cfg, nil
}

func getSecret(ctx context.Context, client *secretmanager.Client, projectID, secretID string) (string, error) {
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID),
	})
	if err != nil {
		return "", fmt.Errorf("access secret %q: %w", secretID, err)
	}
	return strings.TrimSpace(string(result.Payload.Data)), nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
