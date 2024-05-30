package output

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/x-ethr/server/handler/types"
)

type Processor[Output interface{}] func(r *http.Request, output chan<- *types.Response[Output], exception chan<- *types.Exception, options *types.Options)

func Process[Output interface{}](w http.ResponseWriter, r *http.Request, processor Processor[Output], settings ...types.Variadic) {
	ctx := r.Context()

	configuration := types.Configuration()
	for _, option := range settings {
		option(configuration)
	}

	output, exception := Channels[Output]()

	go processor(r.WithContext(ctx), output, exception, configuration)

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

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.Code)

			json.NewEncoder(w).Encode(response.Payload)

			return
		case e := <-exception:
			var err error = e.Source
			if e.Source == nil {
				err = fmt.Errorf("N/A")
			}

			slog.ErrorContext(ctx, "Error While Processing Request", slog.String("error", err.Error()), slog.String("internal-message", e.Log), slog.String("path", r.URL.Path), slog.String("method", r.Method))
			http.Error(w, e.Error(), e.Code)

			return
		}
	}
}
