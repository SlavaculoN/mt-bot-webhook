package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mt-bot/internal/config"
	"mt-bot/internal/infrastructure/repository"
	adminrouter "mt-bot/internal/infrastructure/router"
	tginfra "mt-bot/internal/infrastructure/telegram"
	"mt-bot/internal/transport/telegram"
	"mt-bot/internal/usecase"
	"net"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env: ", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("ошибка загрузки конфигурации: ", err)
	}
	log.Printf("Конфигурация успешно загружена")
	log.Printf("HTTP сервер будет запущен на %s", cfg.HTTPAddr)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal("ошибка создания Telegram Bot Api:", err)
	}

	log.Printf("Бот успешно авторизован: %s", bot.Self.UserName)

	sessionRepo := repository.NewMemorySessionRepository()

	messageIDRepo := telegram.NewMemoryMessageIDRepository()
	adminrouter := adminrouter.NewDirectionAdminRouter(
		cfg.GeneralAdminChatID,
		cfg.ServicePoeticheskiyChatID,
		cfg.ServiceSikeirosaChatID,
	)

	applicationSender := tginfra.NewTelegramApplicationSender(bot)

	applicationUseCase := usecase.NewApplicationUseCase(
		sessionRepo,
		applicationSender,
		adminrouter,
	)

	telegramHandler := telegram.NewTelegramHandler(
		applicationUseCase,
		bot,
		messageIDRepo,
	)

	webhookHandler := telegram.NewWebhookHandler(
		cfg.WebhookSecret,
		telegramHandler,
	)

	mux := http.NewServeMux()

	mux.Handle("/webhook", webhookHandler)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	listener, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %s", err)
	}

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.Serve(listener)
	}()

	log.Print("HTTP Сервер запущен на:", cfg.HTTPAddr)

	if err := setTelegramWebhook(cfg.BotToken, cfg.WebhookURL, cfg.WebhookSecret); err != nil {
		_ = server.Close()
		log.Fatalf("failed to regiter Webhook: %v", err)
	}

	log.Println("Webhook succeed registration")

	if err := <-serverErr; err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed work HTTPserver: %v", err)
	}

}

func setTelegramWebhook(botToken, webhookURL, secretToken string) error {
	payload := struct {
		URL            string   `json:"url"`
		SecretToken    string   `json:"secret_token"`
		AllowedUpdates []string `json:"allowed_updates"`
	}{
		URL:         webhookURL,
		SecretToken: secretToken,
		AllowedUpdates: []string{
			"message",
			"callback_query",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to format JSON: %s", err)
	}

	apiURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/setWebhook", botToken,
	)

	request, err := http.NewRequest(
		http.MethodPost,
		apiURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to create HHTP response %s", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		safeError := strings.ReplaceAll(err.Error(), botToken, "hidden")
		return fmt.Errorf("failed request setWebhook: %s", safeError)
	}

	defer response.Body.Close()

	var result struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
		Result      bool   `json:"result"`
	}

	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result); err != nil {
		return fmt.Errorf("failed to read response telegram %s", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("telegram reqest HTTP %s %s", response.Status, result.Description)
	}

	if !result.OK {
		return fmt.Errorf("telegram cancel setWebhook: %s", err)
	}

	if !result.Result {
		return fmt.Errorf("telegram confirm register Webhhok")
	}

	return nil
}
