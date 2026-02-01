package usecase

import (
	"context"
	"fmt"
)

func (u *usecases) SaveChatUserPair(ctx context.Context, chatID int64, userID string) error {

	const op = "usecase.SaveChatUserPair"

	err := u.chatStorage.SaveChatIDByUserID(ctx, chatID, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.usersStorage.SaveUserIDByChatID(ctx, chatID, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
