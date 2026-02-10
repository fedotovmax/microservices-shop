package ports

import "context"

type UsersStorage interface {
	FindByChatID(ctx context.Context, chatID int64) (string, error)
	Save(ctx context.Context, chatID int64, userID string) error
}
