package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetFromMetadata(ctx context.Context, key string, fallback ...string) []string {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		values := md.Get(key)
		if len(values) > 0 {
			return values
		}
	}

	if len(fallback) > 0 {
		return fallback
	}

	return nil
}
