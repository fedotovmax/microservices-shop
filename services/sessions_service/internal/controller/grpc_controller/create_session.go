package grpccontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *controller) CreateSession(ctx context.Context, req *sessionspb.CreateSessionRequest) (*sessionspb.CreateSessionResponse, error) {

	const op = "grpc_controller.CreateSession"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewPrepareSessionInput(
		req.Uid, req.UserAgent, req.Ip, req.BypassCode, req.DeviceTrustToken,
	)

	err := input.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	newSession, err := c.usecases.CreateSession(ctx, input)

	if err != nil {
		return c.handleCreateSessionError(locale, keys.CreateSessionInternal, err)
	}

	var trustToken *sessionspb.CreatedTrustToken

	if newSession.TrustToken != nil {
		trustToken = &sessionspb.CreatedTrustToken{
			TrustTokenValue:   newSession.TrustToken.DeviceTrustTokenValue,
			TrustTokenExpTime: timestamppb.New(newSession.TrustToken.DeviceTrustTokenExpTime),
		}
	}

	return &sessionspb.CreateSessionResponse{
		Payload: &sessionspb.CreateSessionResponse_SessionCreated{
			SessionCreated: &sessionspb.SessionCreated{
				AccessToken:    newSession.AccessToken,
				RefreshToken:   newSession.RefreshToken,
				AccessExpTime:  timestamppb.New(newSession.AccessExpTime),
				RefreshExpTime: timestamppb.New(newSession.RefreshExpTime),
				TrustToken:     trustToken,
			},
		},
	}, nil

}

func (c *controller) handleCreateSessionError(locale string, fallback string, err error) (*sessionspb.CreateSessionResponse, error) {
	const op = "controller.grpc.handleCreateSessionError"

	l := c.log.With(slog.String("op", op))

	var newDeviceLoginErr *errs.LoginFromNewIPOrDeviceError
	var userInBlacklistErr *errs.UserSessionsInBlacklistError

	switch {

	case errors.Is(err, errs.ErrBadBypassCode):

		msg, i18nerr := i18n.Local.Get(locale, keys.InvalidCode)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &sessionspb.CreateSessionResponse{
			Payload: &sessionspb.CreateSessionResponse_BadBypassCode{
				BadBypassCode: &sessionspb.BadBypassCode{
					Message: msg,
				},
			},
		}, nil

	case errors.As(err, &newDeviceLoginErr):
		msg, i18nerr := i18n.Local.Get(locale, newDeviceLoginErr.ErrCode)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &sessionspb.CreateSessionResponse{
			Payload: &sessionspb.CreateSessionResponse_LoginFromNewDevice{
				LoginFromNewDevice: &sessionspb.LoginFromNewDeviceOrIP{
					Message:       msg,
					CodeExpiresAt: timestamppb.New(newDeviceLoginErr.CodeExpiresAt),
				},
			},
		}, nil

	case errors.As(err, &userInBlacklistErr):

		msg, i18nerr := i18n.Local.Get(locale, userInBlacklistErr.ErrCode)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &sessionspb.CreateSessionResponse{
			Payload: &sessionspb.CreateSessionResponse_UserInBlacklist{
				UserInBlacklist: &sessionspb.UserInBlacklist{
					Message:       msg,
					LinkExpiresAt: timestamppb.New(userInBlacklistErr.LinkExpiresAt),
				},
			},
		}, nil

	default:
		return nil, c.handleError(locale, fallback, err)
	}
}
