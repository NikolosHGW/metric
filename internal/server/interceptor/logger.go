package interceptor

import (
	"context"
	"time"

	"github.com/NikolosHGW/metric/internal/server/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	duration := time.Since(start)
	st, _ := status.FromError(err)

	logger.Log.Sugar().Infoln(
		"method", info.FullMethod,
		"duration", duration,
		"status", st.Code(),
		"error", err,
	)

	return resp, err
}
