package telemetry

import (
	"context"
	"net/http"

	"github.com/x-ethr/server/internal/keystore"
)

type generic struct {
	keystore.Valuer[string]
}

func (generic) Value(ctx context.Context) string {
	return ctx.Value(key).(string)
}

func (generic) Middleware(next http.Handler) http.Handler {
	// var name = text.Title(key.String(), func(o *text.Options) {
	// 	o.Log = true
	// })
	//
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	ctx := r.Context()
	//
	// 	server := ctx.Value(keystore.Keys().Server()).(string)
	// 	service := ctx.Value(keystore.Keys().Service()).(string)
	//
	// 	{
	// 		value := "enabled"
	//
	// 		slog.Log(ctx, logging.Trace, "Middleware", slog.String("name", name), slog.Group("context", slog.String("key", string(key)), slog.Any("value", map[string]string{"enabled": value})))
	//
	// 		ctx = context.WithValue(ctx, key, value)
	// 	}
	//
	// 	servername := fmt.Sprintf("%s-%s", server, service)
	// 	handler := otelhttp.NewHandler(otelhttp.WithRouteTag(r.URL.Path, next), r.Method, otelhttp.WithServerName(servername), otelhttp.WithFilter(func(request *http.Request) (filter bool) {
	// 		if request.URL.Path == "/health" {
	// 			filter = true
	// 		}
	//
	// 		return
	// 	}))
	//
	// 	handler.ServeHTTP(w, r)
	// })

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
