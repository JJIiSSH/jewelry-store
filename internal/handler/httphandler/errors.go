package httphandler

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func bindErrorMessage(err error) string {
	var validationErr validator.ValidationErrors

	if errors.As(err, &validationErr) {
		return "invalid request body"
	}

	return "invalid json request"
}
