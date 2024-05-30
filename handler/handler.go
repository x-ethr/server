package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/x-ethr/server/handler/input"
	"github.com/x-ethr/server/handler/output"
)

func Input[Input interface{}, Output interface{}](w http.ResponseWriter, r *http.Request, v *validator.Validate, processor input.Processor[Input, Output]) {
	input.Process[Input, Output](w, r, v, processor)
}

func Output[Output interface{}](w http.ResponseWriter, r *http.Request, processor output.Processor[Output]) {
	output.Process[Output](w, r, processor)
}
