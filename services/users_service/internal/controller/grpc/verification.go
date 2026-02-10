package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/internal/usecases"
	"google.golang.org/protobuf/types/known/emptypb"
)

type verification struct {
	userspb.UnimplementedVerificationServiceServer
	log                    *slog.Logger
	verifyEmail            *usecases.VerifyEmailUsecase
	sendNewVerifyEmailLink *usecases.SendNewEmailVerifyLinkUsecase
}

func NewVerification(
	log *slog.Logger,
	verifyEmail *usecases.VerifyEmailUsecase,
	sendNewVerifyEmailLink *usecases.SendNewEmailVerifyLinkUsecase,
) *verification {
	return &verification{
		log:                    log,
		verifyEmail:            verifyEmail,
		sendNewVerifyEmailLink: sendNewVerifyEmailLink,
	}
}

func (c *verification) VerifyEmail(ctx context.Context, req *userspb.VerifyEmailRequest) (*userspb.VerifyEmailResponse, error) {

	const op = "controller.grpc.verification.VerifyEmail"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewUUIDInput()
	input.SetUUID(req.GetLink())

	err := input.Validate(locale, "Link")

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	err = c.verifyEmail.Execute(ctx, input.GetUUID())

	if err != nil {
		return c.handleVerifyEmailResponse(locale, err)
	}

	return &userspb.VerifyEmailResponse{
		Payload: &userspb.VerifyEmailResponse_Ok{
			Ok: &userspb.EmailVerifiedSuccess{
				Message: "OK",
			},
		},
	}, nil
}

func (c *verification) handleVerifyEmailResponse(locale string, err error) (
	*userspb.VerifyEmailResponse, error,
) {
	const op = "controller.grpc.verification.handleVerifyEmailResponse"

	l := c.log.With(slog.String("op", op))

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
		return nil, handleError(l, locale, keys.VerifyEmailInternal, err)
	}
}

func (c *verification) SendNewEmailVerifyLink(ctx context.Context, req *userspb.SendNewEmailVerifyLinkRequest) (*emptypb.Empty, error) {
	const op = "controller.grpc.verification.SendNewEmailVerifyLink"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewUUIDInput()
	input.SetUUID(req.GetUserId())

	err := input.Validate(locale, "UserID")

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	err = c.sendNewVerifyEmailLink.Execute(ctx, input.GetUUID(), locale)

	if err != nil {
		return nil, handleError(l, locale, keys.SendNewVerifyEmailInternal, err)
	}

	return &emptypb.Empty{}, nil
}
