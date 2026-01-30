package grpccontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

func (c *controller) VerifyEmail(ctx context.Context, req *userspb.VerifyEmailRequest) (*userspb.VerifyEmailResponse, error) {

	const op = "controller.grpc.FindUserByID"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewUUIDInput()
	input.SetUUID(req.GetLink())

	err := input.Validate(locale, "Link")

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	err = c.usecases.VerifyEmail(ctx, input.GetUUID())

	if err != nil {
		return c.handleVerifyEmailResponse(locale, keys.GetUserInternal, err)
	}

	return &userspb.VerifyEmailResponse{
		Payload: &userspb.VerifyEmailResponse_Ok{
			Ok: &userspb.EmailVerifiedSuccess{
				Message: "OK",
			},
		},
	}, nil

}

func (c *controller) handleVerifyEmailResponse(locale string, fallbackMsg string, err error) (
	*userspb.VerifyEmailResponse, error,
) {
	const op = "controller.grpc.handleVerifyEmailResponse"

	var verifyLinkExpiredErr *errs.VerifyEmailLinkExpiredError

	switch {

	case errors.Is(err, errs.ErrVerifyEmailLinkNotFound):
		return &userspb.VerifyEmailResponse{
			Payload: &userspb.VerifyEmailResponse_NotFound{
				NotFound: &userspb.VerifyEmailLinkNotFound{
					Message: "Verify link not found",
				},
			},
		}, nil

	case errors.As(err, &verifyLinkExpiredErr):
		return &userspb.VerifyEmailResponse{
			Payload: &userspb.VerifyEmailResponse_LinkExpired{
				LinkExpired: &userspb.VerifyEmailLinkExpired{
					Message: "Verify link is expired",
					UserId:  verifyLinkExpiredErr.UID,
				},
			},
		}, nil

	default:
		return nil, c.handleError(locale, keys.VerifyEmailInternal, err)
	}

}
