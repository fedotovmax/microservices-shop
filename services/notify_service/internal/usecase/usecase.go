package usecase

import (
	"context"
	"log/slog"
)

type Storage interface {
	GetChatIDByUserID(ctx context.Context, userID string) (int64, error)
	GetUserIDByChatID(ctx context.Context, chatID int64) (string, error)
	SaveChatIDByUserID(ctx context.Context, chatID int64, userID string) error
	SaveUserIDByChatID(ctx context.Context, chatID int64, userID string) error
}

type usecases struct {
	log     *slog.Logger
	storage Storage
}

func New(log *slog.Logger, storage Storage) *usecases {
	return &usecases{
		log:     log,
		storage: storage,
	}
}
