package handler

// // Options is the configuration structure optionally mutated via the [Variadic] constructor used throughout the package.
// type Options struct {
// 	handler *Handler
// }
//
// // Variadic represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// // nil checks as all implementations first construct an [Options] reference using packaged default(s).
// type Variadic func(o *Options)
//
// type Handler struct {
// 	Writer  http.ResponseWriter
// 	Request *http.Request
// 	Output  chan<- *types.Response
// 	Error   chan<- *types.Exception
// }
//
// // Configuration represents a default constructor.
// func configuration(w http.ResponseWriter, r *http.Request, output chan<- *types.Response, exception chan<- *types.Exception) *Options {
// 	return &Options{
// 		handler: &Handler{
// 			Writer: w, Request: r, Output: output, Error: exception,
// 		},
// 	}
// }
