package types

import (
	"context"
	"errors"
	"net/http"
)

var Null = errors.New("invalid nil pointer")

// Options is the configuration structure optionally mutated via the [Variadic] constructor used throughout the package.
type Options struct {
	CTX *CTX
}

// Variadic represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct an [Options] reference using packaged default(s).
type Variadic func(o *Options)

// Configuration represents a default constructor.
func Configuration(w http.ResponseWriter, r *http.Request, input *interface{}, output chan<- *Response, exception chan<- *Exception) *Options {
	return &Options{
		CTX: &CTX{
			w: w, r: r, input: input, output: output, exception: exception,
		},
	}
}

type CTX struct {
	w         http.ResponseWriter
	r         *http.Request
	input     interface{}
	output    chan<- *Response
	exception chan<- *Exception
}

func (c *CTX) Writer() http.ResponseWriter {
	return c.w
}

// Context updates the underlying request's context.
func (c *CTX) Context(ctx context.Context) {
	c.r = c.r.WithContext(ctx)
}

func (c *CTX) Request() *http.Request {
	return c.r
}

// Input represents an HTTP request's input.
//
// If the input is found to be nil, an error of type [Null] is returned.
func (c *CTX) Input() (interface{}, error) {
	if c.input == nil {
		return nil, Null
	}

	return c.input, nil
}

// Complete is a wrapper around the channel *Response. It's up to consumers to return immediately following a call to [CTX.Complete].
func (c *CTX) Complete(response *Response) {
	c.output <- response
	return
}

// Error is a wrapper around the channel *Exception. It's up to consumers to return immediately following a call to [CTX.Error].
func (c *CTX) Error(exception *Exception) {
	c.exception <- exception
	return
}

func (c *CTX) Channels() (response chan<- *Response, exception chan<- *Exception) {
	return c.output, c.exception
}
