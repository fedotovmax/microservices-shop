package controller

// import (
// 	"context"

// 	"github.com/fedotovmax/microservices-shop/customer-site/internal/openapiclient"
// )

// type refreshCtxKey struct{}

// type tokenType uint8

// const (
// 	TokenTypeAccess  tokenType = 1
// 	TokenTypeRefresh tokenType = 2
// )

// type authToken struct {
// 	Type  tokenType
// 	Value string
// }

// func createAuthContext(
// 	parentCtx context.Context,
// 	tokens ...authToken,
// ) context.Context {

// 	if parentCtx == nil {
// 		parentCtx = context.Background()
// 	}

// 	ctx := parentCtx

// 	var (
// 		accessToken  string
// 		hasAccess    bool
// 		refreshToken string
// 		hasRefresh   bool
// 	)

// 	for _, token := range tokens {
// 		switch token.Type {
// 		case TokenTypeAccess:
// 			accessToken = token.Value
// 			hasAccess = true

// 		case TokenTypeRefresh:
// 			refreshToken = token.Value
// 			hasRefresh = true
// 		}
// 	}

// 	if hasAccess {
// 		ctx = context.WithValue(
// 			ctx,
// 			openapiclient.ContextAPIKeys,
// 			map[string]openapiclient.APIKey{
// 				bearerAuth: {
// 					Key:    accessToken,
// 					Prefix: bearerAuthPrefix,
// 				},
// 			},
// 		)
// 	}

// 	if hasRefresh {
// 		ctx = context.WithValue(ctx, refreshCtxKey{}, refreshToken)
// 	}

// 	return ctx
// }

// func getRefreshTokenFromCtx(ctx context.Context) (string, bool) {
// 	token, ok := ctx.Value(refreshCtxKey{}).(string)

// 	return token, ok
// }
