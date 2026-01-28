package controller

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"net/http"

// 	"github.com/fedotovmax/microservices-shop/customer-site/internal/openapiclient"
// )

// const bearerAuth = "BearerAuth"
// const bearerAuthPrefix = "Bearer"

// const headerAuthorization = "Authorization"

// const userAgentHeader = "User-Agent"

// const xForwadedForHeader = "X-Forwarded-For"

// type authCtxKey struct{}

// type AuthState struct {
// 	Access      string
// 	Refresh     string
// 	UserAgent   string
// 	IP          string
// 	DeviceToken string
// 	Dirty       bool // true, если токены обновились
// }

// func GetAuthState(ctx context.Context) *AuthState {
// 	s, _ := ctx.Value(authCtxKey{}).(*AuthState)
// 	return s
// }

// const accessTokenCookie = "access_token"
// const refreshTokenCookie = "refresh_token"

// func AuthContext(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		state := &AuthState{}

// 		if c, err := r.Cookie(accessTokenCookie); err == nil {
// 			state.Access = c.Value
// 		}
// 		if c, err := r.Cookie(refreshTokenCookie); err == nil {
// 			state.Refresh = c.Value
// 		}

// 		state.UserAgent = r.Header.Get(userAgentHeader)

// 		ip := r.Header.Get(xForwadedForHeader)
// 		if ip == "" {
// 			ip = r.RemoteAddr
// 		}
// 		state.IP = ip

// 		ctx := context.WithValue(r.Context(), authCtxKey{}, state)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// type refreshedKey struct{}

// type AuthTransport struct {
// 	Base          http.RoundTripper
// 	RefreshClient *openapiclient.APIClient
// }

// func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
// 	ctx := req.Context()
// 	state := GetAuthState(ctx)
// 	if state == nil {
// 		return t.base().RoundTrip(req)
// 	}

// 	// первый запрос
// 	req2 := req.Clone(ctx)
// 	if state.Access != "" {
// 		req2.Header.Set(headerAuthorization, fmt.Sprintf("%s %s", bearerAuthPrefix, state.Access))
// 	}

// 	resp, err := t.base().RoundTrip(req2)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// не 401 или нет refresh токена или уже рефрешились
// 	if resp.StatusCode != http.StatusUnauthorized || state.Refresh == "" || ctx.Value(refreshedKey{}) != nil {
// 		return resp, nil
// 	}

// 	resp.Body.Close()

// 	// делаем refresh через отдельный клиент

// 	openapiCtx := context.WithValue(ctx, openapiclient.ContextAPIKeys, map[string]openapiclient.APIKey{
// 		bearerAuth: {
// 			Key:    state.Access,
// 			Prefix: bearerAuthPrefix,
// 		},
// 	})

// 	newTokens, _, err := t.RefreshClient.CustomersAPI.CustomersSessionRefreshSessionPost(openapiCtx).Dto(openapiclient.SessionspbRefreshSessionRequest{
// 		Ip:           state.IP,
// 		RefreshToken: state.Refresh,
// 		UserAgent:    state.UserAgent,
// 	}).Execute()

// 	if err != nil {
// 		return nil, errors.New("refresh failed: " + err.Error())
// 	}

// 	state.Access = newTokens.AccessToken
// 	state.Refresh = newTokens.RefreshToken
// 	state.Dirty = true

// 	ctx = context.WithValue(ctx, refreshedKey{}, true)

// 	// повтор запроса
// 	req3 := req.Clone(ctx)
// 	req3.Header.Set(headerAuthorization, fmt.Sprintf("%s %s", bearerAuthPrefix, state.Access))

// 	return t.base().RoundTrip(req3)
// }

// func (t *AuthTransport) base() http.RoundTripper {
// 	if t.Base != nil {
// 		return t.Base
// 	}
// 	return http.DefaultTransport
// }

// func AuthCookieWriter(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)

// 		state := GetAuthState(r.Context())
// 		if state == nil || !state.Dirty {
// 			return
// 		}

// 		http.SetCookie(w, &http.Cookie{
// 			Name:     accessTokenCookie,
// 			Value:    state.Access,
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   true,
// 		})
// 		http.SetCookie(w, &http.Cookie{
// 			Name:     refreshTokenCookie,
// 			Value:    state.Refresh,
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   true,
// 		})
// 	})
// }

// func InitAPIClient() *openapiclient.APIClient {
// 	refreshCfg := openapiclient.NewConfiguration()
// 	refreshClient := openapiclient.NewAPIClient(refreshCfg)

// 	cfg := openapiclient.NewConfiguration()
// 	cfg.HTTPClient = &http.Client{
// 		Transport: &AuthTransport{
// 			RefreshClient: refreshClient,
// 		},
// 	}

// 	return openapiclient.NewAPIClient(cfg)
// }
