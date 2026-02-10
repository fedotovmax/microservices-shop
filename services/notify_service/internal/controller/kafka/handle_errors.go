package kafka

import (
	"errors"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/pkg/logger"
)

func (k *kafkaController) handleErrors(err error, commit func(), l *slog.Logger) {

	switch {
	case errors.Is(err, ErrInvalidPayloadForEventType):
		l.Error("invalid payload", logger.Err(err))
		commit()
	case errors.Is(err, errs.ErrEventAlreadyHandled):
		l.Error("event was handled", logger.Err(err))
		commit()
	case errors.Is(err, errs.ErrChatIDNotFound):
		l.Error("user are not subscribed for telegram notifications", logger.Err(err))
		commit()
	case errors.Is(err, errs.ErrSendTelegramMessage):
		//no commit)
		l.Error("cannot send message to telegram, will be retry later", logger.Err(err))
	default:
		commit()
	}
}
