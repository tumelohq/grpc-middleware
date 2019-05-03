package grpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcmap "github.com/tumelohq/grpc-middleware/map"
	grpcmask "github.com/tumelohq/grpc-middleware/mask"
	test "github.com/tumelohq/grpc-middleware/testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Example() {
	serverAddress := "127.0.0.1:8900"
	interceptor := grpc.UnaryInterceptor(
		// wrap errors, then mask them so unknown errors are completely identical to unknown ones
		grpcmiddleware.ChainUnaryServer(
			// masks errors
			grpcmask.UnaryServerInterceptor(
				codes.Internal,
				codes.Unknown,
			),
			// wraps errors
			grpcmap.UnaryServerInterceptor(map[codes.Code]codes.Code{
				codes.Unknown: codes.Internal,
			}),
		),
	)
	grpcServer := grpc.NewServer(interceptor)
	s := TestPingService{}
	test.RegisterTestServiceServer(grpcServer, s)
	l, _ := net.Listen("tcp", serverAddress)

	defer l.Close()
	go grpcServer.Serve(l)
	defer grpcServer.GracefulStop()

	conn, _ := grpc.Dial(serverAddress, grpc.WithInsecure())

	c := test.NewTestServiceClient(conn)

	// internal error is masked
	req := &test.Request{
		Code:    int32(codes.Internal),
		Message: "some sensitive info",
	}
	_, err := c.Ping(context.Background(), req)
	fmt.Println(err)

	// wraps unknown error as internal, then masks it
	req = &test.Request{
		Code:    int32(codes.Unknown),
		Message: "unknown data for unknown status",
	}
	_, err = c.Ping(context.Background(), req)
	fmt.Println(err)

	// ignores errors with not found status
	req = &test.Request{
		Code:    int32(codes.NotFound),
		Message: "entity not found",
	}
	_, err = c.Ping(context.Background(), req)
	fmt.Println(err)

	// Output:
	// rpc error: code = Internal desc = Internal
	// rpc error: code = Internal desc = Internal
	// rpc error: code = NotFound desc = entity not found
}

type TestPingService struct {
	T *testing.T
}

func (s TestPingService) Ping(_ context.Context, r *test.Request) (*test.Empty, error) {
	c := codes.Code(r.GetCode())
	return &test.Empty{}, status.Error(c, r.GetMessage())
}
