package grpcmask_test

import (
	"context"
	"fmt"
	"net"

	grpcmask "github.com/tumelohq/grpc-middleware/mask"
	test "github.com/tumelohq/grpc-middleware/testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func ExampleUnaryServerInterceptor() {
	serverAddress := "127.0.0.1:8900"
	interceptor := grpc.UnaryInterceptor(
		// masks the following status codes
		grpcmask.UnaryServerInterceptor(
			codes.Internal,
			codes.Unknown,
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

	req := &test.Request{
		Code:    int32(codes.Internal),
		Message: "some really sensitive data",
	}
	_, err := c.Ping(context.Background(), req)
	fmt.Println(err)

	req = &test.Request{
		Code:    int32(codes.Unknown),
		Message: "some other sensitive data",
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
	// rpc error: code = Internal desc = Internal
	// rpc error: code = Unknown desc = Unknown
	// rpc error: code = NotFound desc = entity not found
}
