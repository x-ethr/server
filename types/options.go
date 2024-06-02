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
func Configuration(w http.ResponseWriter, r *http.Request, input interface{}, output chan<- *Response, redirect chan<- *Redirect, exception chan<- *Exception) *Options {
	return &Options{
		CTX: &CTX{
			w: w, r: r, input: input, output: output, redirect: redirect, exception: exception,
		},
	}
}

// CTX represents a context for handling HTTP requests.
type CTX struct {
	w         http.ResponseWriter
	r         *http.Request
	input     interface{}
	output    chan<- *Response
	redirect  chan<- *Redirect
	exception chan<- *Exception
}

// Writer returns the http.ResponseWriter associated with the CTX object.
// This allows direct access to the underlying response writer to modify the response.
func (c *CTX) Writer() http.ResponseWriter {
	return c.w
}

// Context sets the context for the CTX object by updating the request's context with the provided context.
func (c *CTX) Context(ctx context.Context) {
	c.r = c.r.WithContext(ctx)
}

// Request returns the http.Request associated with the CTX object.
// This allows direct access to the underlying request to access request data and headers.
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

// Redirect is a wrapper around the channel *Redirect. It's up to consumers to return immediately following a call to [CTX.Redirect].
func (c *CTX) Redirect(response *Redirect) {
	c.redirect <- response
	return
}

// Error is a wrapper around the channel *Exception. It's up to consumers to return immediately following a call to [CTX.Error].
func (c *CTX) Error(exception *Exception) {
	c.exception <- exception
	return
}
