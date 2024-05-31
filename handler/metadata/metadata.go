package metadata

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/x-ethr/server/handler"
	"github.com/x-ethr/server/handler/types"
	"github.com/x-ethr/server/middleware"
)

func processor(x *types.CTX) {
	const name = "metadata"

	ctx := x.Request().Context()

	service := middleware.New().Service().Value(ctx)
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(service).Start(ctx, name)

	defer span.End()

	var payload = map[string]interface{}{
		middleware.New().Service().Value(ctx): map[string]interface{}{
			"path":    middleware.New().Path().Value(ctx),
			"service": middleware.New().Service().Value(ctx),
			"version": middleware.New().Version().Value(ctx).Service,
		},
	}

	x.Complete(&types.Response{Code: http.StatusOK, Payload: payload})
	return
}

// Handler returns metadata service-related information.
var Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	handler.Process(w, r, processor)

	return
})
