package grpccanceled_test

import (
	"context"
	"net"
	"strings"
	"testing"

	grpccanceled "github.com/tumelohq/grpc-middleware/canceled"
	test "github.com/tumelohq/grpc-middleware/testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryServerInterceptor(t *testing.T) {
	t.Parallel()
	serverAddress := "127.0.0.1:8900"
	interceptor := grpc.UnaryInterceptor(
		grpccanceled.UnaryServerInterceptor(),
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
		code            codes.Code
		errString       string
		expectedErrCode codes.Code
	}{
		{codes.OK, "", codes.OK},
		{codes.Canceled, "context canceled", codes.Canceled},
		{codes.Internal, "... context canceled", codes.Canceled},
		{codes.Internal, "rpc error: code = Internal desc = this is a legit error", codes.Internal},
	}

	for _, tt := range tts {
		t.Run(tt.errString, func(t *testing.T) {
			req := &test.Request{
				Code:    int32(tt.code),
				Message: tt.errString,
			}
			_, err := c.Ping(context.Background(), req)
			if status.Code(err) != tt.expectedErrCode {
				t.Errorf("want code %s, got %s", tt.expectedErrCode.String(), status.Code(err).String())
			}
			if err == nil && tt.errString == "" {
				return
			}
			if !strings.HasSuffix(err.Error(), tt.errString) {
				t.Errorf("want suffix %s, got %s", tt.errString, err)
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
