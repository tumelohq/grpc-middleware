# grpc-mask

A gRPC interceptor to mask errors.

This piece of middleware intends to stop sensitive information from being returned to calling services. Pass in a list of codes you wish to mask, and they shall all be returned as internal server errors, with no further information.

## Usage 

Have a look at the test, we spin up an entire client and server, and plug the middleware in.

```go
interceptor := grpc.UnaryInterceptor(
    grpcmask.UnaryServerInterceptor(
        codes.Internal,
        codes.Unknown,
    ),
)
```

This example will mask any function that returns an Internal, or Unknown code.