package input

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/x-ethr/server/handler/types"
)

type Processor func(handler *Handler[Input])

func Process[Input interface{}](w http.ResponseWriter, r *http.Request, v *validator.Validate, processor Processor[Input], settings ...Variadic[Input]) {
	ctx := r.Context()

	var input Input

	output, exception, invalid := channels()

	if message, validators, e := types.Validate(ctx, v, r.Body, &input); e != nil {
		invalid <- &types.Invalid{Validators: validators, Message: message, Source: e}
		return
	}

	o := configuration[Input](w, r, &input, output, exception)
	for _, option := range settings {
		option(o)
	}

	o.handler.Request = o.handler.Request.WithContext(ctx)

	go processor(o.handler)

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
