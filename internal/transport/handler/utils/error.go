package utils

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type ErrorDTO struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    int    `json:"code"`
}

func (e *ErrorDTO) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.Code)
	return nil
}

func InternalServerError(err error) *ErrorDTO {
	return &ErrorDTO{
		Message: err.Error(),
		Status:  http.StatusText(http.StatusInternalServerError),
		Code:    http.StatusInternalServerError,
	}
}

func BadRequestError(err error) *ErrorDTO {
	return &ErrorDTO{
		Message: err.Error(),
		Status:  http.StatusText(http.StatusBadRequest),
		Code:    http.StatusBadRequest,
	}
}

func UnprocessableError(msg string) *ErrorDTO {
	return &ErrorDTO{
		Message: msg,
		Status:  http.StatusText(http.StatusUnprocessableEntity),
		Code:    http.StatusUnprocessableEntity,
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errors.Cause(err)):

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		Render(w, r, InternalServerError(err))
		return
	}
}

var _ render.Renderer = (*ErrorDTO)(nil)
