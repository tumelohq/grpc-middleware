package grpccanceled_test

import (
	"context"
	"fmt"
	"net"

	grpccanceled "github.com/tumelohq/grpc-middleware/canceled"
	test "github.com/tumelohq/grpc-middleware/testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func ExampleUnaryServerInterceptor() {
	serverAddress := "127.0.0.1:8900"
	interceptor := grpc.UnaryInterceptor(
		grpccanceled.UnaryServerInterceptor(),
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

	// non-canceled errors are ignored
	req := &test.Request{
		Code:    int32(codes.Internal),
		Message: "an internal error",
	}
	_, err := c.Ping(context.Background(), req)
	fmt.Println(err)

	req = &test.Request{
		Code:    int32(codes.Internal),
		Message: "context canceled",
	}
	_, err = c.Ping(context.Background(), req)
	fmt.Println(err)

	// Output:
	// rpc error: code = Internal desc = an internal error
	// rpc error: code = Canceled desc = rpc error: code = Internal desc = context canceled
}
