package grpcmask

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor implements the UnaryServerInterceptor interface
func UnaryServerInterceptor(cs ...codes.Code) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			errCode := status.Code(err)
			for _, c := range cs {
				if errCode == c {
					err = status.Error(c, c.String())
				}
			}
			return nil, err
		}
		return resp, nil
	}
}
