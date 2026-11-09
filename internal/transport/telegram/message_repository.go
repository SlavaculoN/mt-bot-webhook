package telegram

import (
	"errors"
	"sync"
)

var ErrMessageIDNotFound = errors.New("message is not found")

type MessageIDRepository interface {
	Save(chatID int64, messageID int) error
	Get(chatID int64) (int, error)
	Delete(chatID int64) error
}

type MemoryMessageIDRepository struct {
	mu       sync.RWMutex
	messages map[int64]int
}

func NewMemoryMessageIDRepository() *MemoryMessageIDRepository {
	return &MemoryMessageIDRepository{
		messages: make(map[int64]int),
	}
}

func (r *MemoryMessageIDRepository) Save(chatID int64, messageID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[chatID] = messageID

	return nil
}

func (r *MemoryMessageIDRepository) Get(chatID int64) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	messageID, ok := r.messages[chatID]
	if !ok {
		return 0, ErrMessageIDNotFound
	}

	return messageID, nil
}

func (r *MemoryMessageIDRepository) Delete(chatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.messages, chatID)

	return nil
}
