package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken                  string
	WebhookURL                string
	WebhookSecret             string
	HTTPAddr                  string
	GeneralAdminChatID        int64
	ServicePoeticheskiyChatID int64
	ServiceSikeirosaChatID    int64
}

func parseChatID(envName string) (int64, error) {
	value, exists := os.LookupEnv(envName)
	if !exists {
		return 0, fmt.Errorf("environment variable %s is not set", envName)
	}

	if value == "" {
		return 0, fmt.Errorf("env variable is empty %s", envName)
	}

	chatID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("incorrect value: %s: %w", envName, err)
	}

	if chatID == 0 {
		return 0, fmt.Errorf("env variable %s must be zero", envName)
	}

	return chatID, nil
}

func requiredEnv(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("env variable %s is not set", name)
	}

	value = strings.TrimSpace(value)

	if value == "" {
		return "", fmt.Errorf("env variable %s is empty", name)
	}

	return value, nil
}

func Load() (Config, error) {
	generalAdminChatID, err := parseChatID("GENERAL_ADMIN_CHAT_ID")
	if err != nil {
		return Config{}, err
	}

	servicePoeticheskiyChatID, err := parseChatID("SERVICE_POETICHESKIY_CHAT_ID")
	if err != nil {
		return Config{}, err
	}

	serviceSikeirosaChatID, err := parseChatID("SERVICE_SIKEIROSA_CHAT_ID")
	if err != nil {
		return Config{}, err
	}

	botToken, err := requiredEnv("BOT_TOKEN")
	if err != nil {
		return Config{}, err
	}

	webhookURL, err := requiredEnv("WEBHOOK_URL")
	if err != nil {
		return Config{}, err
	}

	webhookSecret, err := requiredEnv("WEBHOOK_SECRET")
	if err != nil {
		return Config{}, err
	}

	httpAddr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))

	if httpAddr == "" {
		httpAddr = ":8080"
	}

	return Config{
		BotToken:                  botToken,
		WebhookURL:                webhookURL,
		WebhookSecret:             webhookSecret,
		HTTPAddr:                  httpAddr,
		GeneralAdminChatID:        generalAdminChatID,
		ServicePoeticheskiyChatID: servicePoeticheskiyChatID,
		ServiceSikeirosaChatID:    serviceSikeirosaChatID,
	}, nil
}
