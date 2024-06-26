package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Invalid represents an error resulting from failed validation.
//
//   - Message: the validation's string error. If specified, the code must be [http.StatusUnprocessableEntity]
//     or [http.StatusServiceUnavailable]. If the message is "Internal Validation Error",
//     then the validator for the given request input is invalid.
//   - Validators: a map of field names to validation results.
//   - Source: the source error that caused the invalidation.
type Invalid struct {
	// Message represents the validation's string error.
	//
	// 	- If this value is specified, then the code must be [http.StatusUnprocessableEntity] or [http.StatusServiceUnavailable].
	//	- If the message is "Internal Validation Error", then the validator for the given request input is invalid.
	Message    string     `json:"message,omitempty"`
	Validators Validators `json:"validators,omitempty"`
	Source     error      `json:"error,omitempty"` // Source represents the source error
}

// Error returns a string representation of the Exception. If the Exception's Message is empty,
// it returns the standard HTTP status-text for the given code.
func (i *Invalid) Error() string {
	exception := fmt.Errorf("(%d) %s: %s", http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), i.Message).Error()
	if i.Message == "" {
		exception = fmt.Errorf("(%d) %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).Error()
	} else if i.Message == "Internal Validation Error" {
		exception = fmt.Errorf("(%d) %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).Error()
	}

	return exception
}

// Response - If Validators are present, encode them as JSON response with status code 400.
// Otherwise, if Message is present and not equal to "Internal Validation Error",
// respond with status code 422 (Unprocessable Entity).
// Otherwise, respond with status code 503 (Service Unavailable).
func (i *Invalid) Response(w http.ResponseWriter) {
	if i.Validators != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(i.Validators)

		return
	} else if i.Message != "" && i.Message != "Internal Validation Error" {
		http.Error(w, i.Error(), http.StatusUnprocessableEntity)

		return
	}

	http.Error(w, i.Error(), http.StatusServiceUnavailable)
}

// Exception returns a string representation of the Exception. If the Exception's Message is empty,
// it returns the standard HTTP status-text for the given code.
type Exception struct {
	Code    int    `json:"code,omitempty"`    // Code represents an http status-code.
	Message string `json:"message,omitempty"` // Message represents an http status-message

	Log      string                 `json:"log,omitempty"`      // Log represents an internal log message
	Source   error                  `json:"error,omitempty"`    // Source represents the source error
	Metadata map[string]interface{} `json:"metadata,omitempty"` // Metadata represents internal metadata around the error
}

// Error returns a string representation of the Exception. If the Exception's Message is empty,
// it returns the standard HTTP status-text for the given code.
func (e *Exception) Error() string {
	exception := fmt.Errorf("(%d) %s", e.Code, e.Message).Error()
	if e.Message == "" {
		exception = fmt.Sprintf("(%d) %s", e.Code, http.StatusText(e.Code))
	}

	return exception
}
