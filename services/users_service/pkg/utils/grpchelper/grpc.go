package grpchelper

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/fedotovmax/grpcutils/violations"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampToTimePtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime()
	return &t
}

func GetFromMetadata(ctx context.Context, key string, fallback ...string) []string {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		values := md.Get(key)
		if len(values) > 0 {
			return values
		}
	}

	return fallback

}

func ReturnGRPCInternal(l *slog.Logger, msg string, err error) error {

	l.Warn(err.Error())
	st := status.New(codes.Internal, msg)
	return st.Err()
}

func ReturnGRPCBadRequest(l *slog.Logger, msg string, err error) error {

	var ve violations.ValidationErrors
	if errors.As(err, &ve) {
		fieldviolations := ve.ToRPCViolations()

		badRequest := &errdetails.BadRequest{
			FieldViolations: fieldviolations,
		}

		st := status.New(codes.InvalidArgument, msg)

		withDetails, err := st.WithDetails(badRequest)

		if err != nil {
			return st.Err()
		}
		return withDetails.Err()
	}

	l.Warn(err.Error())
	st := status.New(codes.InvalidArgument, msg)
	return st.Err()

}
