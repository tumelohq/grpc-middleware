package grpccanceled

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor implements the UnaryServerInterceptor interface.
// It maps grpc codes from one code to another.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			// check end of string matches context cancelled string (if so - make it codes.Cancelled)
			if strings.HasSuffix(err.Error(), "context canceled") {
				return resp, status.Error(codes.Canceled, err.Error())
			}
			return resp, err
		}
		return resp, err
	}
}
