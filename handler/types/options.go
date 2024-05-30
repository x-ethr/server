package types

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/x-ethr/server/telemetry"
)

// Options is the configuration structure optionally mutated via the [Variadic] constructor used throughout the package.
type Options struct {
	Tracer trace.Tracer

	Name      string // Name represents the telemetry span name
	Service   string // Service represents the telemetry span's resource service
	Version   string // Version represents the telemetry span's resource version
	Workload  string // Workload represents the telemetry span's workload name
	Component string // Component represents the telemetry span's component name
}

func (o *Options) Span(ctx context.Context) (context.Context, trace.Span) {
	return o.Tracer.Start(ctx, o.Name, trace.WithAttributes(attribute.String("workload", o.Workload)), trace.WithAttributes(telemetry.Resources(ctx, o.Service, o.Version).Attributes()...), trace.WithAttributes(attribute.String("component", o.Component)))
}

// Variadic represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct an [Options] reference using packaged default(s).
type Variadic func(o *Options)

// Configuration represents a default constructor.
func Configuration() *Options {
	return &Options{}
}
