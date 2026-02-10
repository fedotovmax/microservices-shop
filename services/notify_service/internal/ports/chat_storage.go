package ports

import "context"

type ChatStorage interface {
	FindByUID(ctx context.Context, uid string) (int64, error)
	Save(ctx context.Context, chatID int64, userID string) error
}
