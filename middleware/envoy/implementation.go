package envoy

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/x-ethr/text"

	"github.com/x-ethr/server/internal/keystore"
	"github.com/x-ethr/server/logging"
)

type generic struct {
	keystore.Valuer[string]
}

type Envoy struct {
	Attempts *int    `json:"x-envoy-attempt-count"`
	Original *string `json:"x-envoy-original-path"`
	Internal *bool   `json:"x-envoy-internal"`
}

func (*generic) Value(ctx context.Context) *Envoy {
	return ctx.Value(key).(*Envoy)
}

func (*generic) Middleware(next http.Handler) http.Handler {
	var name = text.Title(key.String(), func(o *text.Options) {
		o.Log = true
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		{
			var internal *bool
			if v := r.Header.Get("X-Envoy-Internal"); v == "true" {
				assignment := true
				internal = &assignment
			}

			var attempts *int
			if v := r.Header.Get("X-Envoy-Request-Count"); v != "" {
				assignment, e := strconv.Atoi(v)
				if e == nil {
					attempts = &assignment
				}
			}

			var original *string
			if v := r.Header.Get("X-Envoy-Original-Path"); v != "" {
				original = &v
			}

			value := &Envoy{Original: original, Attempts: attempts, Internal: internal}

			slog.Log(ctx, logging.Trace, "Middleware", slog.String("name", name), slog.Group("context", slog.String("key", string(key)), slog.Any("value", value)))

			ctx = context.WithValue(ctx, key, value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
