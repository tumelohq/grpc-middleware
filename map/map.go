package grpcmap

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor implements the UnaryServerInterceptor interface.
// It maps grpc codes from one code to another.
func UnaryServerInterceptor(cm map[codes.Code]codes.Code) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			errCode := status.Code(err)
			if code, ok := cm[errCode]; ok {
				return nil, status.Error(code, status.Convert(err).Message())
			}
			return nil, err
		}
		return resp, nil
	}
}
