package interceptor

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(logger *zap.Logger) grpc.UnaryServerInterceptor {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			logger.Error("panic", zap.Any("stack", p))

			return status.Errorf(codes.Internal, ErrInternalServer)
		}),
	}

	return recovery.UnaryServerInterceptor(recoveryOpts...)
}
