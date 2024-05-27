# `server` - HTTP Routing, Logging & Telemetry

## Task-Board

- [ ] Complete Documentation
- [ ] Middleware Unit-Testing

## Documentation

Official `godoc` documentation (with examples) can be found at the [Package Registry](https://pkg.go.dev/github.com/x-ethr/server).

## Usage

###### Add Package Dependency

```bash
go get -u github.com/x-ethr/server
```

###### Import & Implement

`main.go`

```go
package main

import (
    "context"
    "encoding/json"
    "errors"
    "flag"
    "fmt"
    "log/slog"
    "net"
    "net/http"
    "os"
    "path/filepath"
    "runtime"

    "github.com/x-ethr/server"
    "github.com/x-ethr/server/logging"
    "github.com/x-ethr/server/middleware"
    "github.com/x-ethr/server/middleware/name"
    "github.com/x-ethr/server/middleware/servername"
    "github.com/x-ethr/server/middleware/versioning"
    "github.com/x-ethr/server/telemetry"
)

// header is a dynamically linked string value - defaults to "server" - which represents the server name.
var header string = "server"

// service is a dynamically linked string value - defaults to "service" - which represents the service name.
var service string = "service"

// version is a dynamically linked string value - defaults to "development" - which represents the service's version.
var version string = "development" // production builds have version dynamically linked

// ctx, cancel represent the server's runtime context and cancellation handler.
var ctx, cancel = context.WithCancel(context.Background())

// port represents a cli flag that sets the server listening port
var port = flag.String("port", "8080", "Server Listening Port.")

var (
    tracer = otel.Tracer(service)
)

func main() {
    // Create Handler Instance
    mux := server.New()

    // Setup Middleware(s)

    // --> Initialize Path Middleware
    mux.Middleware(middleware.New().Path().Middleware)

    // --> Initialize Timeout Middleware
    mux.Middleware(middleware.New().Timeout().Configuration(func(options *timeout.Settings) {
        options.Timeout = time.Second * 3
    }).Middleware)

    // --> Initialize Server-Name Middleware
    mux.Middleware(middleware.New().Server().Configuration(func(options *servername.Settings) {
        options.Server = header
    }).Middleware)

    // --> Initialize Service-Name Middleware
    mux.Middleware(middleware.New().Service().Configuration(func(options *name.Settings) {
        options.Service = service
    }).Middleware)

    // --> Initialize Versioning Middleware
    mux.Middleware(middleware.New().Version().Configuration(func(options *versioning.Settings) {
        options.Version.API = os.Getenv("VERSION")
        if options.Version.API == "" && os.Getenv("CI") == "" {
            options.Version.API = "local"
        }

        options.Version.Service = version
    }).Middleware)

    // --> Initialize Telemetry Middleware
    mux.Middleware(middleware.New().Telemetry().Middleware)

    // Establish Route(s)
    mux.Register("GET /", func(w http.ResponseWriter, r *http.Request) {
        ctx, span := tracer.Start(r.Context(), "example")
        defer span.End()

        path := middleware.New().Path().Value(ctx)

        local := ctx.Value(http.LocalAddrContextKey).(net.Addr)
        output := map[string]interface{}{
            "path":  path,
            "local": local,
        }

        span.SetAttributes(attribute.String("path", path))

        span.SetStatus(codes.Ok, "successful http response")

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(output)
        return
    })

    // --> Example Timeout Usage
    mux.Register("GET /timeout", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        process := time.Duration(rand.Intn(5)) * time.Second

        select {
        case <-ctx.Done():
            return

        case <-time.After(process):
            // The above channel simulates some hard work.
        }

        w.Write([]byte("done"))
    })

    // Start the HTTP Server
    slog.Info("Starting Server ...", slog.String("local", fmt.Sprintf("http://localhost:%s", *(port))))

    // --> Initialize Server Instance
    api := server.Server(ctx, mux, *port)

    // --> Issue Cancellation Handler
    server.Interrupt(ctx, cancel, api)

    // --> Setup Telemetry
    shutdown, e := telemetry.Setup(ctx, service, version)
    if e != nil {
        panic(e)
    }

    defer func() {
        e = errors.Join(e, shutdown(ctx))
    }()

    // <-- Blocking: Begin the Server Listener
    if e := api.ListenAndServe(); e != nil && !(errors.Is(e, http.ErrServerClosed)) {
        slog.ErrorContext(ctx, "Error During Server's Listen & Serve Call ...", slog.String("error", e.Error()))

        os.Exit(100)
    }

    // --> Exit: End the Program Upon Signal
    {
        slog.InfoContext(ctx, "Graceful Shutdown Complete")

        // Waiter
        <-ctx.Done()
    }
}

func init() {
    flag.Parse()

    level := slog.Level(-8)
    if os.Getenv("CI") == "true" {
        level = slog.LevelDebug
    }

    logging.Level(level)

    if service == "service" && os.Getenv("CI") != "true" {
        _, file, _, ok := runtime.Caller(0)
        if ok {
            service = filepath.Base(filepath.Dir(file))
        }
    }

    handler := logging.Logger(func(o *logging.Options) { o.Service = service })
    logger := slog.New(handler)
    slog.SetDefault(logger)
}
```

- Please refer to the [code examples](./examples) for additional usage and implementation details.
- See https://pkg.go.dev/github.com/x-ethr/server for additional documentation.

## Contributions

See the [**Contributing Guide**](./CONTRIBUTING.md) for additional details on getting started.
