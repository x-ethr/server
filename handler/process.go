package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Processor[Input interface{}, Output interface{}] func(ctx context.Context, input *Input, output chan<- *Output, exception chan<- *Exception)

func Process[Input interface{}, Output interface{}](w http.ResponseWriter, r *http.Request, v *validator.Validate, processor Processor[Input, Output]) {
	ctx := r.Context()

	var input Input // only used for logging

	output, exception, invalid := Channels[Output]()

	if message, validators, e := Validate(ctx, v, r.Body, &input); e != nil {
		invalid <- &Invalid{Validators: validators, Message: message, Source: e}
		return
	}

	go processor(ctx, &input, output, exception)

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

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			json.NewEncoder(w).Encode(response)

			return
		case e := <-invalid:
			slog.WarnContext(ctx, "Invalid Request", slog.String("error", e.Error()), slog.String("path", r.URL.Path), slog.String("method", r.Method), slog.Any("input", input))
			e.Response(w)
			return
		case e := <-exception:
			var err error = e.Source
			if e.Source == nil {
				err = fmt.Errorf("N/A")
			}

			slog.ErrorContext(ctx, "Error While Processing Request", slog.String("error", err.Error()), slog.String("internal-message", e.Log), slog.String("path", r.URL.Path), slog.String("method", r.Method), slog.Any("input", input))
			http.Error(w, e.Error(), e.Code)

			return
		}
	}
}
