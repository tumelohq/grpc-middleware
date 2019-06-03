package grpcctxextractor

import (
	"context"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor implements the UnaryServerInterceptor interface
// Interface is a map with string key which is the key for the log. Plus a function value
// which extracts the value from context and potentially returning an error.
func UnaryServerInterceptor(in map[string]func(context.Context) (string, error)) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tags := grpc_ctxtags.Extract(ctx)
		for k, f := range in {
			s, err := f(ctx)
			if err != nil {
				return nil, err
			}
			tags = tags.Set(k, s)
		}
		return handler(ctx, req)
	}
}
