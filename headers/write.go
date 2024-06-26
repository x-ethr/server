package headers

import "net/http"

// Write ...
//
// Headers include:
//   - "portal"
//   - "device"
//   - "user"
//   - "travel"
//   - "x-request-id"
//   - "x-b3-traceid"
//   - "x-b3-spanid"
//   - "x-b3-parentspanid"
//   - "x-b3-sampled"
//   - "x-b3-flags"
//   - "x-ot-span-context"
func Write(r *http.Response, writer http.ResponseWriter) {
	headers := []string{
		http.CanonicalHeaderKey("portal"),
		http.CanonicalHeaderKey("device"),
		http.CanonicalHeaderKey("user"),
		http.CanonicalHeaderKey("travel"),
		http.CanonicalHeaderKey("x-request-id"),
		http.CanonicalHeaderKey("x-b3-traceid"),
		http.CanonicalHeaderKey("x-b3-spanid"),
		http.CanonicalHeaderKey("x-b3-parentspanid"),
		http.CanonicalHeaderKey("x-b3-sampled"),
		http.CanonicalHeaderKey("x-b3-flags"),
		http.CanonicalHeaderKey("x-ot-span-context"),
		http.CanonicalHeaderKey("x-api-version"),
	}

	for key := range headers {
		header := headers[key]

		assignment := r.Header.Get(header)

		if assignment != "" {
			writer.Header().Set(header, assignment)
		}
	}
}
