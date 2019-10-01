package grpcmap_test

import (
	"context"
	"net"
	"reflect"
	"testing"

	grpcmap "github.com/tumelohq/grpc-middleware/map"
	test "github.com/tumelohq/grpc-middleware/testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryServerInterceptor(t *testing.T) {
	t.Parallel()
	serverAddress := "127.0.0.1:8901"
	interceptor := grpc.UnaryInterceptor(
		grpcmap.UnaryServerInterceptor(map[codes.Code]codes.Code{
			codes.Unknown: codes.Internal,
		}),
	)
	grpcServer := grpc.NewServer(interceptor)
	s := TestPingService{}
	test.RegisterTestServiceServer(grpcServer, s)
	l, err := net.Listen("tcp", serverAddress)
	if err != nil {
		t.Fatalf("can't listen to %s: %v", serverAddress, err)
	}
	defer l.Close()
	go grpcServer.Serve(l)
	defer grpcServer.GracefulStop()

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("can't dial %s: %v", serverAddress, err)
	}
	defer conn.Close()
	c := test.NewTestServiceClient(conn)

	tts := []struct {
		code     codes.Code
		expected error
	}{
		{codes.OK, nil},
		{codes.Unknown, status.Error(codes.Internal, "Unknown")},
		{codes.Internal, status.Error(codes.Internal, "Internal")},
		{codes.NotFound, status.Error(codes.NotFound, "NotFound")},
	}

	for _, tt := range tts {
		t.Run(tt.code.String(), func(t *testing.T) {
			req := &test.Request{
				Code:    int32(tt.code),
				Message: tt.code.String(),
			}
			_, err := c.Ping(context.Background(), req)
			if !reflect.DeepEqual(tt.expected, err) {
				t.Errorf("want %s, got %s", tt.expected, err)
			}
		})
	}
}

type TestPingService struct {
	T *testing.T
}

func (s TestPingService) Ping(_ context.Context, r *test.Request) (*test.Empty, error) {
	c := codes.Code(r.GetCode())
	return &test.Empty{}, status.Error(c, r.GetMessage())
}
