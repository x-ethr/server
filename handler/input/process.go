package input

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/x-ethr/server/handler/types"
)

type Processor[Input interface{}] func(w http.ResponseWriter, r *http.Request, input *Input, output chan<- *types.Response, exception chan<- *types.Exception, options *types.Options)

func Process[Input interface{}](w http.ResponseWriter, r *http.Request, v *validator.Validate, processor Processor[Input], settings ...types.Variadic) {
	ctx := r.Context()

	configuration := types.Configuration()
	for _, option := range settings {
		option(configuration)
	}

	var input Input // only used for logging

	output, exception, invalid := channels()

	if message, validators, e := types.Validate(ctx, v, r.Body, &input); e != nil {
		invalid <- &types.Invalid{Validators: validators, Message: message, Source: e}
		return
	}

	go processor(w, r.WithContext(ctx), &input, output, exception, configuration)

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
					w.Write([]byte(http.StatusText(http.StatusNoContent)))
					return
				}

				w.Write([]byte(response.Payload.(string)))
				return
			default:
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response.Payload)
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
