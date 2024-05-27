package name_test

import (
	"net/http"

	"github.com/x-ethr/server/middleware"
	"github.com/x-ethr/server/middleware/name"
)

func Example() {
	mux := http.NewServeMux()

	handler := middleware.New().Service().Configuration(func(options *name.Settings) {
		options.Service = "Service-Name-Value"
	}).Middleware(mux)

	http.ListenAndServe(":8080", handler)
}
