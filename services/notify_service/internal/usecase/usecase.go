package usecase

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
)

type Storage interface {
	GetChatIDByUserID(ctx context.Context, userID string) (int64, error)
	GetUserIDByChatID(ctx context.Context, chatID int64) (string, error)
	SaveChatIDByUserID(ctx context.Context, chatID int64, userID string) error
	SaveUserIDByChatID(ctx context.Context, chatID int64, userID string) error
}

type TgSender interface {
	SendMessage(ctx context.Context, n *domain.TgNotification) error
}

type usecases struct {
	log      *slog.Logger
	storage  Storage
	tgSender TgSender
}

func New(log *slog.Logger, storage Storage, tgSender TgSender) *usecases {
	return &usecases{
		log:      log,
		storage:  storage,
		tgSender: tgSender,
	}
}
