package telemetry

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Labeler(ctx context.Context) *otelhttp.Labeler {
	labeler, _ := otelhttp.LabelerFromContext(ctx)

	return labeler
}
