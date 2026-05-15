package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

func readRequest[T any](r *http.Request) (T, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, err
	}

	return v, nil
}

func writeResponse(w http.ResponseWriter, r *http.Request, v any, status int) {
	withRequestLogContextStatus(r, status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		withRequestLogContextError(r, err)
	}
}

func writeBadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	response := BadRequestErrorResponse
	status := http.StatusBadRequest

	withRequestLogContextStatus(r, status)
	withRequestLogContextError(r, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		withRequestLogContextError(r, err)
	}
}

func writeDomainOrInternalError(w http.ResponseWriter, r *http.Request, err error) {
	response, status := errorToDomainOrInternalErrorResponse(err)

	withRequestLogContextStatus(r, status)
	withRequestLogContextError(r, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		withRequestLogContextError(r, err)
	}
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

var InternalServerErrorResponse = ErrorResponse{
	Error: InternalServerError,
}

var BadRequestErrorResponse = ErrorResponse{
	Error: BadRequestError,
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var InternalServerError = Error{
	Code:    "internal",
	Message: "internal server error",
}

var BadRequestError = Error{
	Code:    "invalid_request",
	Message: "invalid request",
}

func errorToDomainOrInternalErrorResponse(err error) (ErrorResponse, int) {
	domainErr, ok := errors.AsType[*domain.Error](err)
	if !ok {
		return InternalServerErrorResponse, http.StatusInternalServerError
	}

	return domainErrorToErrorResponse(domainErr)
}

func domainErrorToErrorResponse(domainErr *domain.Error) (ErrorResponse, int) {
	var status int

	switch {
	case errors.Is(domainErr, domain.ErrDepartmentNotFound):
		status = http.StatusNotFound
	case errors.Is(domainErr, domain.ErrDepartmentParentNotFound):
		status = http.StatusNotFound
	case errors.Is(domainErr, domain.ErrEmployeeDepartmentNotFound):
		status = http.StatusNotFound
	case errors.Is(domainErr, domain.ErrDepartmentParentConflict):
		status = http.StatusConflict
	case errors.Is(domainErr, domain.ErrDepartmentNameConflict):
		status = http.StatusConflict
	default:
		status = http.StatusBadRequest
	}

	return ErrorResponse{
		Error: Error{
			Code:    domainErr.Code,
			Message: domainErr.Message,
		},
	}, status
}
