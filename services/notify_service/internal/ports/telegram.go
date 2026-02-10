package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
)

type TelegramSender interface {
	SendMessage(ctx context.Context, n *domain.TgNotification) error
}
