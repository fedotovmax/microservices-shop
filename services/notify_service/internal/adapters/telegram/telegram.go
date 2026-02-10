package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
)

type telegram struct {
	tgbot *bot.Bot
	ctx   context.Context
	stop  context.CancelFunc
}

type Config struct {
	Token   string
	Options []bot.Option
}

func New(cfg *Config) (*telegram, error) {

	const op = "adapters.telegram.New"

	tgbot, err := bot.New(cfg.Token, cfg.Options...)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	setCmdCtx, cancelSetCmdCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelSetCmdCtx()

	err = setCommands(setCmdCtx, tgbot)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &telegram{
		tgbot: tgbot,
		ctx:   ctx,
		stop:  cancel,
	}, nil

}

func (tg *telegram) Start() {
	go tg.tgbot.Start(tg.ctx)
}

func (tg *telegram) Stop() {
	tg.stop()

}
