package input

import (
	"net/http"

	"github.com/x-ethr/server/handler/types"
)

// Options is the configuration structure optionally mutated via the [Variadic] constructor used throughout the package.
type Options[Input interface{}] struct {
	handler *Handler[Input]
}

// Variadic represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct an [Options] reference using packaged default(s).
type Variadic[Input interface{}] func(o *Options[Input])

type Handler[Input interface{}] struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Input   *Input
	Output  chan<- *types.Response
	Error   chan<- *types.Exception
}

// configuration represents a default constructor.
func configuration[Input interface{}](w http.ResponseWriter, r *http.Request, input *Input, output chan<- *types.Response, exception chan<- *types.Exception) *Options[Input] {
	return &Options[Input]{
		handler: &Handler[Input]{
			Writer: w, Request: r, Output: output, Error: exception, Input: input,
		},
	}
}
