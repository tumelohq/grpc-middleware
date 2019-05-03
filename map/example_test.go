package grpcmap_test

import (
	"context"
	"fmt"
	"net"

	grpcmap "github.com/tumelohq/grpc-middleware/map"
	test "github.com/tumelohq/grpc-middleware/testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func ExampleUnaryServerInterceptor() {
	serverAddress := "127.0.0.1:8900"
	interceptor := grpc.UnaryInterceptor(
		grpcmap.UnaryServerInterceptor(map[codes.Code]codes.Code{
			codes.Unknown: codes.Internal,
		}),
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

	req := &test.Request{
		Code:    int32(codes.Internal),
		Message: "an internal error",
	}
	_, err := c.Ping(context.Background(), req)
	fmt.Println(err)

	req = &test.Request{
		Code:    int32(codes.Unknown),
		Message: "an unknown error",
	}
	_, err = c.Ping(context.Background(), req)
	fmt.Println(err)

	req = &test.Request{
		Code:    int32(codes.NotFound),
		Message: "entity not found",
	}
	_, err = c.Ping(context.Background(), req)
	fmt.Println(err)

	// Output:
	// rpc error: code = Internal desc = an internal error
	// rpc error: code = Internal desc = an unknown error
	// rpc error: code = NotFound desc = entity not found
}
