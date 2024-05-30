package tracing

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/x-ethr/text"
	"go.opentelemetry.io/otel/trace"

	"github.com/x-ethr/server/internal/keystore"
	"github.com/x-ethr/server/logging"
)

type generic struct {
	keystore.Valuer[string]

	options *Settings
}

func (g *generic) Configuration(options ...Variadic) Implementation {
	var o = settings()
	for _, option := range options {
		option(o)
	}

	g.options = o

	return g
}

func (*generic) Value(ctx context.Context) trace.Tracer {
	return ctx.Value(key).(trace.Tracer)
}

func (g *generic) Middleware(next http.Handler) http.Handler {
	var name = text.Title(key.String(), func(o *text.Options) {
		o.Log = true
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		{
			value := g.options.Tracer

			slog.Log(ctx, logging.Trace, "Middleware", slog.String("name", name), slog.Group("context", slog.String("key", string(key))))

			ctx = context.WithValue(ctx, key, value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
