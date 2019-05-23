package grpclogtracing

import (
	"context"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor implements the UnaryServerInterceptor interface
// Appends trace_id and span_id of the request to the set of tags.
// To be used with grpc_ctxtags and grpc_logging.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tags := grpc_ctxtags.Extract(ctx)
		tags = tags.Set("span.trace_id", trace.FromContext(ctx).SpanContext().TraceID.String())
		tags = tags.Set("span.span_id", trace.FromContext(ctx).SpanContext().SpanID.String())
		return handler(ctx, req)
	}
}
