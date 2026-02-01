package usecase

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
)

type UsersStorage interface {
	SaveUserIDByChatID(ctx context.Context, chatID int64, userID string) error
	GetUserIDByChatID(ctx context.Context, chatID int64) (string, error)
}

type ChatStorage interface {
	GetChatIDByUserID(ctx context.Context, userID string) (int64, error)
	SaveChatIDByUserID(ctx context.Context, chatID int64, userID string) error
}

type EventStorage interface {
	FindEvent(ctx context.Context, eventID string) (string, error)
	SaveEventID(ctx context.Context, eventID string) error
}

type TgSender interface {
	SendMessage(ctx context.Context, n *domain.TgNotification) error
}

type usecases struct {
	log           *slog.Logger
	usersStorage  UsersStorage
	chatStorage   ChatStorage
	eventsStorage EventStorage
	tgSender      TgSender
}

func New(
	log *slog.Logger,
	usersStorage UsersStorage,
	chatStorage ChatStorage,
	eventsStorage EventStorage,
	tgSender TgSender,
) *usecases {
	return &usecases{
		log:           log,
		usersStorage:  usersStorage,
		chatStorage:   chatStorage,
		eventsStorage: eventsStorage,
		tgSender:      tgSender,
	}
}
