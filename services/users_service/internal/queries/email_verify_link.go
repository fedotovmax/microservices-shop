package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/ports"
)

type EmailVerifyLink interface {
	FindByUserID(ctx context.Context, uid string) (*domain.EmailVerifyLink, error)

	Find(ctx context.Context, link string) (*domain.EmailVerifyLink, error)
}

type emailVerifyLink struct {
	emailVerifyStorage ports.EmailVerifyStorage
}

func NewEmailVerifyLink(emailVerifyStorage ports.EmailVerifyStorage) EmailVerifyLink {
	return &emailVerifyLink{
		emailVerifyStorage: emailVerifyStorage,
	}
}

func (q *emailVerifyLink) FindByUserID(
	ctx context.Context,
	uid string,
) (*domain.EmailVerifyLink, error) {

	linkEntity, err := q.emailVerifyStorage.FindBy(ctx, db.VerifyEmailLinkUserIDField, uid)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrVerifyEmailLinkNotFound, err)
		}
		return nil, err
	}

	return linkEntity, nil
}

func (q *emailVerifyLink) Find(
	ctx context.Context,
	link string,
) (*domain.EmailVerifyLink, error) {

	linkEntity, err := q.emailVerifyStorage.FindBy(ctx, db.VerifyEmailLinkPrimaryField, link)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrVerifyEmailLinkNotFound, err)
		}
		return nil, err
	}

	return linkEntity, nil
}
