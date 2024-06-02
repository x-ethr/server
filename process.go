package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/x-ethr/server/logging"
	"github.com/x-ethr/server/types"
)

// Handle represents a function type that takes a pointer to [types.CTX] as an argument.
//
//   - The function is used to handle the processing of an HTTP request and produce a response.
//   - The function can be used as a parameter in functions like [Validate] and [Process].
//
// Example: @TODO write real example
//
//	    func MyHandle(x *types.CTX) {
//			   // handle the request and produce a response
//	    }
//	    // usage in Validate function
//	    Validate(w http.ResponseWriter, r *http.Request, v *validator.Validate, handle Handle, settings ...types.Variadic)
//
//	    // usage in Process function
//	    Process(w http.ResponseWriter, r *http.Request, handle Handle, settings ...types.Variadic)
type Handle func(x *types.CTX)

// Validate is an enhanced version of [Process]. Specifically, with a validator.Validate as an argument, users of [Processor] will be able
// to use [types.CTX] [types.CTX.Input] function to retrieve a hydrated instance of the input data structure from the request's body.
//
//   - The [types.CTX.Input] value will be a pointer to the [Input] Generic specification from the caller.
func Validate[Input interface{}](w http.ResponseWriter, r *http.Request, v *validator.Validate, handle Handle, settings ...types.Variadic) {
	ctx := r.Context()

	output, redirect, exception := channels()

	invalid := make(chan *types.Invalid)

	var input Input
	if message, validators, e := types.Validate(ctx, v, r.Body, &input); e != nil {
		invalid <- &types.Invalid{Validators: validators, Message: message, Source: e}
	}

	o := types.Configuration(w, r, nil, output, redirect, exception)
	for _, option := range settings {
		option(o)
	}

	o.CTX.Context(ctx)

	go handle(o.CTX)

	for {
		select {
		case <-ctx.Done():
			return
		case response := <-output:
			evaluate(ctx, w, r, response)
			return
		case response := <-redirect:
			redirection(ctx, w, r, response)
			return
		case e := <-invalid:
			slog.WarnContext(ctx, "Invalid Request", slog.String("error", e.Error()), slog.String("path", r.URL.Path), slog.String("method", r.Method), slog.Any("input", input))
			e.Response(w)
			return
		case e := <-exception:
			throw(ctx, w, r, e)
			return
		}
	}
}

// Process is a function that handles the processing of an HTTP request. It takes an http.ResponseWriter, an *http.Request,
// a Processor function, and optional settings as parameters. The function initializes the necessary channels,
// creates an options object, and starts a goroutine to execute the processor function. It then enters a loop,
// waiting for either a response or an exception. If a response is received, it writes the response to the writer
// and returns. If an exception is received, it logs the error and returns an HTTP error with the corresponding status code.
func Process(w http.ResponseWriter, r *http.Request, handle Handle, settings ...types.Variadic) {
	ctx := r.Context()

	output, redirect, exception := channels()

	o := types.Configuration(w, r, nil, output, redirect, exception)
	for _, option := range settings {
		option(o)
	}

	o.CTX.Context(ctx)

	go handle(o.CTX)

	for {
		select {
		case <-ctx.Done():
			return
		case response := <-output:
			evaluate(ctx, w, r, response)
			return
		case response := <-redirect:
			redirection(ctx, w, r, response)
			return
		case e := <-exception:
			throw(ctx, w, r, e)
			return
		}
	}
}

func throw(ctx context.Context, w http.ResponseWriter, request *http.Request, e *types.Exception) {
	var err error = e.Source
	if e.Source == nil {
		err = fmt.Errorf("N/A")
	}

	slog.ErrorContext(ctx, "Error While Processing Request", slog.Any("metadata", e.Metadata), slog.String("error", err.Error()), slog.String("public", e.Message), slog.String("internal", e.Log), slog.String("path", request.URL.Path), slog.String("method", request.Method))
	http.Error(w, e.Error(), e.Code)

	return
}

func redirection(ctx context.Context, w http.ResponseWriter, request *http.Request, redirect *types.Redirect) {
	if redirect == nil {
		slog.ErrorContext(ctx, "Redirect Returned Unexpected, Null Result", slog.String("path", request.URL.Path), slog.String("method", request.Method))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Log(ctx, logging.Trace, "Successfully Processed Redirect", slog.Any("redirect", redirect))

	http.Redirect(w, request, redirect.URL, redirect.Status)
	return
}

func evaluate(ctx context.Context, w http.ResponseWriter, request *http.Request, response *types.Response) {
	if response == nil {
		slog.ErrorContext(ctx, "Response Returned Unexpected, Null Result", slog.String("path", request.URL.Path), slog.String("method", request.Method))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Log(ctx, logging.Trace, "Successfully Processed Request", slog.Any("response", response))

	w.WriteHeader(response.Status)

	switch response.Payload.(type) {
	case []byte:
		value := response.Payload.([]byte)
		w.Header().Set("Content-Type", "text/plain")

		if value == nil {
			w.WriteHeader(http.StatusNoContent)
			slog.WarnContext(ctx, "No Content from Response", slog.Any("response", response))
			if size, e := w.Write([]byte(http.StatusText(http.StatusNoContent))); e != nil {
				slog.ErrorContext(ctx, "Unable to Write Response Body (Text, No-Content)", slog.String("error", e.Error()), slog.Int("size", size))
			}

			return
		}

		if size, e := w.Write(value); e != nil {
			slog.ErrorContext(ctx, "Unable to Write Response Body (Text)", slog.String("error", e.Error()), slog.Int("size", size))
		}

		return
	case string:
		value := response.Payload.(string)
		w.Header().Set("Content-Type", "text/plain")

		if value == "" {
			w.WriteHeader(http.StatusNoContent)
			slog.WarnContext(ctx, "No Content from Response", slog.Any("response", response))
			if size, e := w.Write([]byte(http.StatusText(http.StatusNoContent))); e != nil {
				slog.ErrorContext(ctx, "Unable to Write Response Body (Text, No-Content)", slog.String("error", e.Error()), slog.Int("size", size))
			}

			return
		}

		if size, e := w.Write([]byte(value)); e != nil {
			slog.ErrorContext(ctx, "Unable to Write Response Body (Text)", slog.String("error", e.Error()), slog.Int("size", size))
		}

		return
	case *string:
		value := response.Payload.(*string)
		w.Header().Set("Content-Type", "text/plain")

		if value == nil {
			w.WriteHeader(http.StatusNoContent)
			slog.WarnContext(ctx, "No Content from Response", slog.Any("response", response))
			if size, e := w.Write([]byte(http.StatusText(http.StatusNoContent))); e != nil {
				slog.ErrorContext(ctx, "Unable to Write Response Body (Text, No-Content)", slog.String("error", e.Error()), slog.Int("size", size))
			}

			return
		}

		if size, e := w.Write([]byte(*value)); e != nil {
			slog.ErrorContext(ctx, "Unable to Write Response Body (Text)", slog.String("error", e.Error()), slog.Int("size", size))
		}

		return
	default:
		w.Header().Set("Content-Type", "application/json")
		if e := json.NewEncoder(w).Encode(response.Payload); e != nil {
			slog.ErrorContext(ctx, "Unable to Write Response Body (JSON)", slog.String("error", e.Error()))
		}

		return
	}
}
