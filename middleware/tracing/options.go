package tracing

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/x-ethr/server/internal/keystore"
)

type Settings struct {
	Tracer trace.Tracer
}

type Variadic keystore.Variadic[Settings]

func settings() *Settings {
	return &Settings{}
}
