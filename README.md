# grpc-mask

[![GoDoc](https://godoc.org/github.com/tumelohq/grpc-middleware?status.svg)](https://godoc.org/github.com/tumelohq/grpc-middleware)

A gRPC interceptor to mask errors.

This piece of middleware intends to stop sensitive information from being returned to calling services. Pass in a list of codes you wish to mask, and they shall all be returned as internal server errors, with no further information.

## Go get

```sh
go get github.com/tumelohq/grpc-middleware/...
```

## Usage 

Have a look at the test, we spin up an entire client and server, and plug the middleware in.

### Masking

This masks all internal and unknown errors, hiding potentially damaging and sensitive information.

```go
interceptor := grpc.UnaryInterceptor(
    grpcmask.UnaryServerInterceptor(
        codes.Internal,
        codes.Unknown,
    ),
)
```

### Mapping

This maps from one error to another. In this snippet we are mapping Unknown errors to Internal.

```go
interceptor := grpc.UnaryInterceptor(
    grpcmap.UnaryServerInterceptor(map[codes.Code]codes.Code{
        codes.Unknown: codes.Internal,
    }),
)
```
