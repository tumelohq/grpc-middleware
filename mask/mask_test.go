package grpcmask_test

import (
	"context"
	"net"
	"reflect"
	"testing"

	grpcmask "github.com/tumelohq/grpc-middleware/mask"
	test "github.com/tumelohq/grpc-middleware/testing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryServerInterceptor(t *testing.T) {
	t.Parallel()
	serverAddress := "127.0.0.1:8900"
	interceptor := grpc.UnaryInterceptor(
		grpcmask.UnaryServerInterceptor(
			codes.Internal,
			codes.Unknown,
		),
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
		{codes.Canceled, status.Error(codes.Canceled, "Canceled")},
		{codes.Unknown, status.Error(codes.Unknown, "Unknown")},
		{codes.InvalidArgument, status.Error(codes.InvalidArgument, "InvalidArgument")},
		{codes.DeadlineExceeded, status.Error(codes.DeadlineExceeded, "DeadlineExceeded")},
		{codes.NotFound, status.Error(codes.NotFound, "NotFound")},
		{codes.AlreadyExists, status.Error(codes.AlreadyExists, "AlreadyExists")},
		{codes.PermissionDenied, status.Error(codes.PermissionDenied, "PermissionDenied")},
		{codes.ResourceExhausted, status.Error(codes.ResourceExhausted, "ResourceExhausted")},
		{codes.FailedPrecondition, status.Error(codes.FailedPrecondition, "FailedPrecondition")},
		{codes.Aborted, status.Error(codes.Aborted, "Aborted")},
		{codes.OutOfRange, status.Error(codes.OutOfRange, "OutOfRange")},
		{codes.Unimplemented, status.Error(codes.Unimplemented, "Unimplemented")},
		{codes.Internal, status.Error(codes.Internal, "Internal")},
		{codes.Unavailable, status.Error(codes.Unavailable, "Unavailable")},
		{codes.DataLoss, status.Error(codes.DataLoss, "DataLoss")},
		{codes.Unauthenticated, status.Error(codes.Unauthenticated, "Unauthenticated")},
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
