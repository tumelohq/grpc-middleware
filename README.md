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

### Mask

This masks all internal and unknown errors, hiding potentially damaging and sensitive information.

```go
interceptor := grpc.UnaryInterceptor(
    grpcmask.UnaryServerInterceptor(
        codes.Internal,
        codes.Unknown,
    ),
)
```

### Map

This maps from one error to another. In this snippet we are mapping Unknown errors to Internal.

```go
interceptor := grpc.UnaryInterceptor(
    grpcmap.UnaryServerInterceptor(map[codes.Code]codes.Code{
        codes.Unknown: codes.Internal,
    }),
)
```

### Logtracing

Used with [tags](`github.com/grpc-ecosystem/go-grpc-middleware/tags`) and [logging](`github.com/grpc-ecosystem/go-grpc-middleware/logging`), allows you to log out a trace and span ID for each call. Logging middleware logs metadata about the response of the call, such as status code, error message, time, duration, and service name. Using tags and this middleware allow us to add span and trace id to each request, allowing us to corrolate a trace with particular logs (and maybe users too with the ctx extractor below)

```go
// where 'l' is a log entry
interceptor := grpc.UnaryInterceptor(
    grpc.UnaryInterceptor(
        grpcMiddleware.ChainUnaryServer(
            grpcTags.UnaryServerInterceptor(),
            grpcLogTracing.UnaryServerInterceptor(),
            grpcLogging.UnaryServerInterceptor(l),
        ),
    ),
)
```

### ctx extractor

Also used with [logging](`github.com/grpc-ecosystem/go-grpc-middleware/logging`), allows you to pass in a function that extracts a string from context. Typically used for functional requirements, like a user ID, or whatever else may be kept in context.

```go
// where 'l' is a log entry
// and 'enricher' is a function that extracts a string from context
interceptor := grpc.UnaryInterceptor(
    grpc.UnaryInterceptor(
        grpcMiddleware.ChainUnaryServer(
            grpcExtractor.UnaryServerInterceptor(enricher),
            grpcLogging.UnaryServerInterceptor(l),
        ),
    ),
)
```

### Canceled

Used to capture canceled requests. When propagating errors through microservices, we may find canceled errors are badly wrapped and get returned to the user as an internal error. With heavy monitoring, this may trigger alerts for a high % of internal errors for false alarms. This middleware takes canceled errors and ensures they are returned with the correct error code.

```go
interceptor := grpc.UnaryInterceptor(
    grpccanceled.UnaryServerInterceptor(),
)
```
