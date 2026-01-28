package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/microservices-shop/customer-site/internal/keys"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/openapiclient"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/state"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/logger"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/utils"
	"github.com/fedotovmax/singlecall"
)

var ErrUnauthorized = errors.New("unauthorized")

func ClearCookies(w http.ResponseWriter) {
	// access_token
	http.SetCookie(w, &http.Cookie{
		Name:     keys.CookieAccessToken,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // если у тебя HTTPS
		MaxAge:   -1,   // удаляем cookie
	})

	// refresh_token
	http.SetCookie(w, &http.Cookie{
		Name:     keys.CookieRefreshToken,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	//TODO: add other cookies
}

func CreateOpenapiAccessCtx(ctx context.Context, token string) context.Context {
	openapiCtx := context.WithValue(ctx, openapiclient.ContextAPIKeys, map[string]openapiclient.APIKey{
		keys.BearerAuth: {
			Key:    token,
			Prefix: keys.BearerAuthPrefix,
		},
	})

	return openapiCtx
}

type Call[T any] func(ctx context.Context) (T, *http.Response, error)

type Refresher func(ctx context.Context, dto openapiclient.SessionspbRefreshSessionRequest) (*openapiclient.SessionspbSessionCreated, *http.Response, error)

type RefreshResponse struct {
	NewSession   *openapiclient.SessionspbSessionCreated
	HttpResponse *http.Response
}

var refreshGroup singlecall.Group[string, *RefreshResponse]

const singleCallRefreshTokenKey = "RefreshToken"

func WithAuth[T any](
	ctx context.Context,
	log *slog.Logger,
	w http.ResponseWriter,
	state *state.ClientState,
	refresher Refresher,
	call Call[T],
) (T, *http.Response, error) {

	var zero T

	res, httpres, err := call(ctx)

	if err == nil {
		return res, httpres, nil
	}

	if httpres.StatusCode != http.StatusUnauthorized {
		return res, httpres, err
	}

	refreshResonse, _, err := refreshGroup.Do(
		ctx,
		fmt.Sprintf("%s:%s", singleCallRefreshTokenKey, state.RefreshToken),
		func(ctx context.Context) (*RefreshResponse, error) {
			newTokens, refreshHttpRes, err := refresher(ctx, openapiclient.SessionspbRefreshSessionRequest{
				Ip:           state.IP,
				RefreshToken: state.RefreshToken,
				UserAgent:    state.UserAgent,
			})

			return &RefreshResponse{
				NewSession:   newTokens,
				HttpResponse: refreshHttpRes,
			}, err
		})

	if err != nil {
		log.Error("Error when refreshing session", logger.Err(err))
		return zero, refreshResonse.HttpResponse, ErrUnauthorized
	}

	state.AccessToken = refreshResonse.NewSession.AccessToken
	state.RefreshToken = refreshResonse.NewSession.RefreshToken

	http.SetCookie(w, &http.Cookie{
		Name:  keys.CookieAccessToken,
		Value: refreshResonse.NewSession.AccessToken,
		Path:  "/",
		Expires: utils.TimestamppbToTime(
			refreshResonse.NewSession.AccessExpTime.Seconds,
			refreshResonse.NewSession.AccessExpTime.Nanos,
		),
		HttpOnly: true,
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  keys.CookieRefreshToken,
		Value: refreshResonse.NewSession.RefreshToken,
		Path:  "/",
		Expires: utils.TimestamppbToTime(
			refreshResonse.NewSession.RefreshExpTime.Seconds,
			refreshResonse.NewSession.RefreshExpTime.Nanos,
		),
		HttpOnly: true,
		Secure:   true,
	})

	ctx = CreateOpenapiAccessCtx(ctx, state.AccessToken)

	return call(ctx)

}
