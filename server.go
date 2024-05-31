package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/x-ethr/server/internal/writer"
)

// Server initializes a http.Server with application-specific configuration.
func Server(ctx context.Context, handler http.Handler, middleware *Middlewares, port string) *http.Server {
	handler = writer.Handle(middleware.Handler(handler))

	if v, ok := ctx.Value("server-name").(string); ok {
		handler = otelhttp.NewHandler(handler, "server", otelhttp.WithServerName(v), otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))
	} else {
		handler = otelhttp.NewHandler(handler, "server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))
	}

	return &http.Server{
		Addr:                         fmt.Sprintf("0.0.0.0:%s", port),
		Handler:                      handler,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  15 * time.Second,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 60 * time.Second,
		IdleTimeout:                  30 * time.Second,
		MaxHeaderBytes:               http.DefaultMaxHeaderBytes,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		ConnContext: nil,
	}
}
