package telegram

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/usecases"
	"github.com/go-telegram/bot"
)

type Telegram interface {
	RegisterCommand(cmdType bot.HandlerType, cmd domain.Cmd, f bot.HandlerFunc) error
}

type tgBotController struct {
	log              *slog.Logger
	saveChatUserPair *usecases.SaveChatUserPairUsecase
	tg               Telegram
}

// https://t.me/MicroservicesShopNotifyBot?start=12345
func New(log *slog.Logger, saveChatUserPair *usecases.SaveChatUserPairUsecase, tg Telegram) *tgBotController {
	return &tgBotController{
		log:              log,
		saveChatUserPair: saveChatUserPair,
		tg:               tg,
	}
}

func (tgc *tgBotController) Register() error {

	const op = "controller.tg_bot.Register"

	err := tgc.tg.RegisterCommand(bot.HandlerTypeMessageText, domain.Start, tgc.handleStart)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tgc.tg.RegisterCommand(bot.HandlerTypeMessageText, domain.Help, tgc.handleHelp)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

// returning cmd:string, args:[]string, error
func (tgc *tgBotController) parseCommand(text string) (string, []string, error) {

	parts := strings.Fields(text)

	if len(parts) == 0 {
		return "", nil, errs.ErrInvalidCommand
	}

	cmd := tgc.normalizeCmd(parts[0])

	var args []string

	if len(parts) > 1 {
		args = parts[1:]
	}

	return cmd, args, nil

}

func (tgc *tgBotController) normalizeCmd(cmd string) string {

	if i := strings.Index(cmd, "@"); i != -1 {
		return cmd[:i]
	}
	return cmd
}
