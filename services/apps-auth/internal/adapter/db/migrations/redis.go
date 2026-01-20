package migrations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter/db/redisadapter"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/utils"
)

func ApplyRedisMigrations(ctx context.Context, adminSecret string, rdb *redisadapter.Rdb) error {
	const op = "adapter.db.migrations.ApplyRedisMigrations"

	secretHash := utils.CreateHash(adminSecret)

	_, err := rdb.FindApp(ctx, secretHash)

	if err != nil && !errors.Is(err, adapter.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return nil
	}

	err = rdb.SaveApp(ctx, secretHash, &domain.App{CreatedAt: time.Now().UTC(), Name: "MAIN_ADMIN_APP", Type: domain.ApplicationTypeAdmin})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
