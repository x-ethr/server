package types

import (
	"go.opentelemetry.io/otel/trace"
)

// Options is the configuration structure optionally mutated via the [Variadic] constructor used throughout the package.
type Options struct {
	Tracer  trace.Tracer
	Service string
	Version string
}

// Variadic represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct an [Options] reference using packaged default(s).
type Variadic func(o *Options)

// Configuration represents a default constructor.
func Configuration() *Options {
	return &Options{}
}
