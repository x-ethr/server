package middleware

import (
	"github.com/x-ethr/server/middleware/envoy"
	"github.com/x-ethr/server/middleware/name"
	"github.com/x-ethr/server/middleware/path"
	"github.com/x-ethr/server/middleware/servername"
	"github.com/x-ethr/server/middleware/state"
	"github.com/x-ethr/server/middleware/telemetry"
	"github.com/x-ethr/server/middleware/timeout"
	"github.com/x-ethr/server/middleware/tracing"
	"github.com/x-ethr/server/middleware/versioning"
)

type generic struct{}

func (*generic) Path() path.Implementation {
	return path.New()
}

func (*generic) Version() versioning.Implementation {
	return versioning.New()
}

func (*generic) Service() name.Implementation {
	return name.New()
}

func (*generic) Telemetry() telemetry.Implementation {
	return telemetry.New()
}

func (*generic) Server() servername.Implementation {
	return servername.New()
}

func (*generic) Timeout() timeout.Implementation {
	return timeout.New()
}

func (*generic) Envoy() envoy.Implementation {
	return envoy.New()
}

func (*generic) Tracer() tracing.Implementation {
	return tracing.New()
}

func (*generic) State() state.Implementation {
	return state.New()
}

type Middleware interface {
	Path() path.Implementation           // Path - See the [path] package for additional details.
	Version() versioning.Implementation  // Version - See the [versioning] package for additional details.
	Service() name.Implementation        // Service - See the [name] package for additional details.
	Telemetry() telemetry.Implementation // Telemetry - See the [telemetry] package for additional details.
	Server() servername.Implementation   // Server - See the [servername] package for additional details.
	Timeout() timeout.Implementation     // Timeout - See the [timeout] package for additional details.
	Envoy() envoy.Implementation         // Envoy - See the [envoy] package for additional details.
	Tracer() tracing.Implementation      // Tracer - See the [tracing] package for additional details.
	State() state.Implementation         // State - See the [state] package for additional details.
}

func New() Middleware {
	return &generic{}
}
