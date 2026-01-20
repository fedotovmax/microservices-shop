package interceptors

import (
	"context"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/appsauthpb"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/keys"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AdminSecretInterceptor(secret string) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		if info.FullMethod == appsauthpb.AppsAuthService_CreateApp_FullMethodName {
			md, ok := metadata.FromIncomingContext(ctx)

			if !ok {
				return nil, status.Error(codes.Unauthenticated, "missing metadata")
			}

			adminKeys := md.Get(keys.MetadataAdminSecret)
			expected := "123"
			if len(adminKeys) == 0 || adminKeys[0] == expected {
				return nil, status.Error(codes.Unauthenticated, "invalid admin secret key")
			}
		}

		return handler(ctx, req)
	}
}
