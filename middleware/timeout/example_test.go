package timeout_test

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/x-ethr/server/middleware"
	"github.com/x-ethr/server/middleware/timeout"
)

func Example() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		process := time.Duration(rand.Intn(5)) * time.Second

		select {
		case <-ctx.Done():
			return

		case <-time.After(process):
			// The above channel simulates some hard work.
		}

		w.Write([]byte("done"))
	})

	handler := middleware.New().Timeout().Configuration(func(options *timeout.Settings) {
		options.Timeout = time.Second * 3
	}).Middleware(mux)

	http.ListenAndServe(":8080", handler)
}
