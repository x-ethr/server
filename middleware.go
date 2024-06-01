package server

import (
	"net/http"
)

func Middleware() *Middlewares {
	return &Middlewares{
		middleware: make([]func(http.Handler) http.Handler, 0),
	}
}

type Middlewares struct {
	middleware []func(http.Handler) http.Handler
}

func (m *Middlewares) Add(middlewares ...func(http.Handler) http.Handler) {
	if len(middlewares) == 0 {
		return
	}

	m.middleware = append(m.middleware, middlewares...)
}

func (m *Middlewares) Handler(parent http.Handler) (handler http.Handler) {
	var length = len(m.middleware)
	if length == 0 {
		return parent
	}

	// Wrap the end handler with the middleware chain
	handler = m.middleware[len(m.middleware)-1](parent)
	for i := len(m.middleware) - 2; i >= 0; i-- {
		handler = m.middleware[i](handler)
	}

	return
}
