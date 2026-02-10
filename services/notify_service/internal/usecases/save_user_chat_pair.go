package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/ports"
)

type SaveChatUserPairUsecase struct {
	log          *slog.Logger
	chatStorage  ports.ChatStorage
	usersStorage ports.UsersStorage
}

func NewSaveChatUserPairUsecase(
	log *slog.Logger,
	chatStorage ports.ChatStorage,
	usersStorage ports.UsersStorage,
) *SaveChatUserPairUsecase {
	return &SaveChatUserPairUsecase{
		log:          log,
		chatStorage:  chatStorage,
		usersStorage: usersStorage,
	}
}

func (u *SaveChatUserPairUsecase) Execute(ctx context.Context, chatID int64, userID string) error {

	const op = "usecases.save_user_chat_pair"

	err := u.chatStorage.Save(ctx, chatID, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.usersStorage.Save(ctx, chatID, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
