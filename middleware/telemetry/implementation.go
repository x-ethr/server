package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/x-ethr/text"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/x-ethr/server/internal/keystore"
	"github.com/x-ethr/server/logging"
)

type generic struct {
	keystore.Valuer[string]
}

func (generic) Value(ctx context.Context) string {
	return ctx.Value(key).(string)
}

func (generic) Middleware(next http.Handler) http.Handler {
	var name = text.Title(key.String(), func(o *text.Options) {
		o.Log = true
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		server := ctx.Value(keystore.Keys().Server()).(string)
		service := ctx.Value(keystore.Keys().Service()).(string)

		{
			value := "enabled"

			slog.Log(ctx, logging.Trace, "Middleware", slog.String("name", name), slog.Group("context", slog.String("key", string(key)), slog.Any("value", map[string]string{"enabled": value})))

			ctx = context.WithValue(ctx, key, value)
		}

		servername := fmt.Sprintf("%s-%s", server, service)
		handler := otelhttp.NewHandler(otelhttp.WithRouteTag(r.URL.Path, next), r.Method, otelhttp.WithServerName(servername), otelhttp.WithFilter(func(request *http.Request) (filter bool) {
			ctx := request.Context()

			if request.URL.Path == "/health" {
				filter = true

				slog.Log(ctx, logging.Debug, "Health Telemetry Exclusion", slog.Bool("filter", filter))
			}

			return
		}))

		handler.ServeHTTP(w, r)
	})
}
