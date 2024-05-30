package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type Validators map[string]Validator

type Validator struct {
	Value   interface{} `json:"value,omitempty"`
	Valid   bool        `json:"valid"`
	Message string      `json:"message"`
}

type Helper interface {
	Help() Validators
}

func Validate(ctx context.Context, v *validator.Validate, body io.Reader, data Helper) (string, Validators, error) {
	// invalid describes an invalid argument passed to `Struct`, `StructExcept`, StructPartial` or `Field`
	var invalid *validator.InvalidValidationError

	// Unmarshal request-body into "data".
	if e := json.NewDecoder(body).Decode(&data); e != nil {
		// Log an issue unmarshalling the body and return a Bad request exception.
		slog.Log(ctx, slog.LevelWarn, "Unable to Unmarshal Request Body",
			slog.String("error", e.Error()),
		)

		return "Valid JSON Required as Input", nil, e
	}

	// Validate "data" using "validation".
	if e := v.Struct(data); e != nil {
		// Check if the error is due to an invalid validation configuration.
		if errors.As(e, &invalid) {
			// Log the issue and return an Internal server error exception.
			slog.ErrorContext(ctx, "Invalid Validator", slog.String("error", e.Error()))

			return "Internal Validation Error", nil, e
		}

		// Loop through the validation errors, logging each one.
		for key, e := range e.(validator.ValidationErrors) {
			slog.Log(ctx, slog.LevelWarn, fmt.Sprintf("Validator (%d)", key), slog.Group("error",
				slog.String("tag", e.Tag()),
				slog.String("actual-tag", e.ActualTag()),
				slog.String("parameter", e.Param()),
				slog.String("field", e.Field()),
				slog.String("namespace", e.Namespace()),
				slog.String("struct-namespace", e.StructNamespace()),
				slog.String("struct-field", e.StructField()),
				slog.Group("reflection",
					slog.Attr{
						Key:   "kind",
						Value: slog.AnyValue(e.Kind()),
					},
					slog.Attr{
						Key:   "type",
						Value: slog.AnyValue(e.Type()),
					},
					slog.Attr{
						Key:   "value",
						Value: slog.AnyValue(e.Value()),
					},
				),
			))
		}

		return "", data.Help(), e
	}

	// If no exception was generated, log the "data" for debugging purposes.
	// slog.Log(ctx, slog.LevelDebug, "Request", logging.Structure("body", data))

	// Return nil if there were no exceptions generated.
	return "", nil, nil
}
