package telegram

import (
	"context"
	"encoding/json"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateProcessor interface {
	HandleUpdate(ctx context.Context, update tgbotapi.Update) error
}

type WebhookHandler struct {
	secret    string
	processor UpdateProcessor
}

func NewWebhookHandler(secret string, processor UpdateProcessor) *WebhookHandler {
	return &WebhookHandler{
		secret:    secret,
		processor: processor,
	}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	secret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")

	if secret != h.secret {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var update tgbotapi.Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.processor.HandleUpdate(r.Context(), update); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
