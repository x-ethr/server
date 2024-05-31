package output

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/x-ethr/server/handler/types"
)

type Processor func(c *types.CTX)

// Validate is an enhanced version of [Process]. Specifically, with a validator.Validate as an argument, users of [Processor] will be able
// to use [types.CTX] [types.CTX.Input] function to retrieve a hydrated instance of the input data structure from the request's body.
func Validate(w http.ResponseWriter, r *http.Request, v *validator.Validate, processor Processor, settings ...types.Variadic) {
	ctx := r.Context()

	output, exception := channels()

	invalid := make(chan *types.Invalid)

	var input interface{}
	if message, validators, e := types.Validate(ctx, v, r.Body, &input); e != nil {
		invalid <- &types.Invalid{Validators: validators, Message: message, Source: e}
		return
	}

	o := types.Configuration(w, r, input, output, exception)
	for _, option := range settings {
		option(o)
	}

	o.CTX.Context(ctx)

	go processor(o.CTX)

	for {
		select {
		case <-ctx.Done():
			return
		case response := <-output:
			if response == nil {
				slog.ErrorContext(ctx, "Response Returned Unexpected, Null Result", slog.String("path", r.URL.Path), slog.String("method", r.Method), slog.Any("input", input))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			slog.DebugContext(ctx, "Successfully Processed Request", slog.Any("response", response))

			w.WriteHeader(response.Code)

			switch response.Payload.(type) {
			case string, *string:
				w.Header().Set("Content-Type", "text/plain")
				if response.Payload == nil {
					if size, e := w.Write([]byte(http.StatusText(http.StatusNoContent))); e != nil {
						slog.ErrorContext(ctx, "Unable to Write Response Body (Text)", slog.String("error", e.Error()), slog.Int("size", size))
					}

					return
				}

				if size, e := w.Write([]byte(http.StatusText(http.StatusNoContent))); e != nil {
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
		case e := <-invalid:
			slog.WarnContext(ctx, "Invalid Request", slog.String("error", e.Error()), slog.String("path", r.URL.Path), slog.String("method", r.Method), slog.Any("input", input))
			e.Response(w)
			return
		case e := <-exception:
			var err error = e.Source
			if e.Source == nil {
				err = fmt.Errorf("N/A")
			}

			slog.ErrorContext(ctx, "Error While Processing Request", slog.Any("metadata", e.Metadata), slog.String("error", err.Error()), slog.String("public", e.Message), slog.String("internal", e.Log), slog.String("path", r.URL.Path), slog.String("method", r.Method))
			http.Error(w, e.Error(), e.Code)

			return
		}
	}
}

// Process is a function that handles the processing of an HTTP request. It takes an http.ResponseWriter, an *http.Request,
// a Processor function, and optional settings as parameters. The function initializes the necessary channels,
// creates an options object, and starts a goroutine to execute the processor function. It then enters a loop,
// waiting for either a response or an exception. If a response is received, it writes the response to the writer
// and returns. If an exception is received, it logs the error and returns an HTTP error with the corresponding status code.
func Process(w http.ResponseWriter, r *http.Request, processor Processor, settings ...types.Variadic) {
	ctx := r.Context()

	output, exception := channels()

	o := types.Configuration(w, r, nil, output, exception)
	for _, option := range settings {
		option(o)
	}

	o.CTX.Context(ctx)

	go processor(o.CTX)

	for {
		select {
		case <-ctx.Done():
			return
		case response := <-output:
			if response == nil {
				slog.ErrorContext(ctx, "Response Returned Unexpected, Null Result", slog.String("path", r.URL.Path), slog.String("method", r.Method))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			slog.DebugContext(ctx, "Successfully Processed Request", slog.Any("response", response))

			w.WriteHeader(response.Code)

			switch response.Payload.(type) {
			case string, *string:
				w.Header().Set("Content-Type", "text/plain")
				if response.Payload == nil {
					if size, e := w.Write([]byte(http.StatusText(http.StatusNoContent))); e != nil {
						slog.ErrorContext(ctx, "Unable to Write Response Body (Text)", slog.String("error", e.Error()), slog.Int("size", size))
					}

					return
				}

				if size, e := w.Write([]byte(response.Payload.(string))); e != nil {
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
		case e := <-exception:
			var err error = e.Source
			if e.Source == nil {
				err = fmt.Errorf("N/A")
			}

			slog.ErrorContext(ctx, "Error While Processing Request", slog.Any("metadata", e.Metadata), slog.String("error", err.Error()), slog.String("public", e.Message), slog.String("internal", e.Log), slog.String("path", r.URL.Path), slog.String("method", r.Method))
			http.Error(w, e.Error(), e.Code)

			return
		}
	}
}