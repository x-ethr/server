package servername_test

import (
	"net/http"

	"github.com/x-ethr/server/middleware"
	"github.com/x-ethr/server/middleware/servername"
)

func Example() {
	mux := http.NewServeMux()

	handler := middleware.New().Server().Configuration(func(options *servername.Settings) {
		options.Server = "Server-Name-Value"
	}).Middleware(mux)

	http.ListenAndServe(":8080", handler)
}
