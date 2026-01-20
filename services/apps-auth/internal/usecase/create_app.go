package usecase

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/utils"
)

func (u *usecases) CreateApp(ctx context.Context, newApp *domain.App) (string, error) {
	const op = "usecases.CreateApp"

	secret, err := utils.CreateSecret()

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	hash := utils.CreateHash(secret)

	err = u.storage.SaveApp(ctx, hash, newApp)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return secret, nil

}
