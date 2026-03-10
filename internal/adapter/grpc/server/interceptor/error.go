package interceptor

import (
	"context"
	"errors"

	"github.com/go-list-templ/grpc/internal/core/domain/entityerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrDeadlineExceeded = "request server timeout"
	ErrInternalServer   = "internal server"
)

func ErrorHandling() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, toGrpcError(err)
		}

		return resp, nil
	}
}

func toGrpcError(err error) error {
	switch {
	case errors.Is(err, entityerr.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entityerr.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entityerr.ErrUserInvalidData):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, ErrDeadlineExceeded)
	default:
		return status.Error(codes.Internal, ErrInternalServer)
	}
}
