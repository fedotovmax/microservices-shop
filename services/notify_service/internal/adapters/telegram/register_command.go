package telegram

import (
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/go-telegram/bot"
)

func (tg *telegram) RegisterCommand(cmdType bot.HandlerType, cmd domain.Cmd, f bot.HandlerFunc) error {

	const op = "adapters.telegram.RegisterCommand"

	err := cmd.Validate()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tg.tgbot.RegisterHandler(cmdType, cmd.String(), bot.MatchTypePrefix, f)
	return nil
}
