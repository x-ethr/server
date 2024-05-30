package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Processor[Body interface{}, Output interface{}] func(ctx context.Context, reader io.ReadCloser, body chan *Body, output chan<- *Output, exception chan<- *Exception, invalid chan<- *Invalid)

func Process[Body interface{}, Output interface{}](w http.ResponseWriter, r *http.Request, processor Processor[Body, Output]) {
	ctx := r.Context()

	var input *Body // only used for logging

	body, output, exception, invalid := Channels[Body, Output]()

	go processor(ctx, r.Body, body, output, exception, invalid)

	for {
		select {
		case <-ctx.Done():
			return
		case input = <-body: // continue waiting for one of the other primitives to complete
			continue
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
