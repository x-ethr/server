package types

import "net/http"

// Options is the configuration structure optionally mutated via the [Variadic] constructor used throughout the package.
type Options struct {
	handler *Handler
}

// Variadic represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct an [Options] reference using packaged default(s).
type Variadic func(o *Options)

// Configuration represents a default constructor.
func Configuration(w http.ResponseWriter, r *http.Request, output chan<- *Response, exception chan<- *Exception) *Options {
	return &Options{
		handler: &Handler{
			w: w, r: r, output: output, exception: exception,
		},
	}
}

type Handler struct {
	w         http.ResponseWriter
	r         *http.Request
	output    chan<- *Response
	exception chan<- *Exception
	options   *Options
}
