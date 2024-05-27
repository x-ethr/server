package path_test

import (
	"net/http"

	"github.com/x-ethr/server/middleware"
)

func Example() {
	mux := http.NewServeMux()

	handler := middleware.New().Path().Middleware(mux)

	http.ListenAndServe(":8080", handler)
}
