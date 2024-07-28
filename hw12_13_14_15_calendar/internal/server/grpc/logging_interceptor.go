package grpc

import (
	"context"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Infof("gRPC method: %s, request: %v", info.FullMethod, req)
		resp, err := handler(ctx, req)
		if err != nil {
			st, _ := status.FromError(err)
			logger.Errorf("gRPC method: %s, error: %v, code: %v", info.FullMethod, err, st.Code())
		} else {
			logger.Infof("gRPC method: %s, response: %v", info.FullMethod, resp)
		}
		return resp, err
	}
}
