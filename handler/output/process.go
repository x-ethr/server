package output

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Processor func(handler *Handler)

func Process(w http.ResponseWriter, r *http.Request, processor Processor, settings ...Variadic) {
	ctx := r.Context()

	output, exception := channels()

	// ctx = context.WithValue(ctx, "writer", w)
	// ctx = context.WithValue(ctx, "request", r)
	// ctx = context.WithValue(ctx, "options", configuration)

	o := configuration(w, r, output, exception)
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
