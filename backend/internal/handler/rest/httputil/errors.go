package httputil

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"errors"
	"net/http"
)

type ErrorMapping struct {
	Err    error
	Status int
	Code   string
}

func WriteServiceError(w http.ResponseWriter, err error, mappings ...ErrorMapping) {
	for _, m := range mappings {
		if errors.Is(err, m.Err) {
			response.SendErrorResponse(w, m.Status, m.Code, err, err.Error())
			return
		}
	}
	if status, ok := customErr.HTTPStatusFromError(err); ok {
		response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
		return
	}
	response.SendErrorResponse(w, http.StatusInternalServerError, "internal_server_error", err, "internal server error")
}

func WriteServiceErrorOr(w http.ResponseWriter, err error, fallbackCode, fallbackMsg string, mappings ...ErrorMapping) {
	for _, m := range mappings {
		if errors.Is(err, m.Err) {
			response.SendErrorResponse(w, m.Status, m.Code, err, err.Error())
			return
		}
	}
	if status, ok := customErr.HTTPStatusFromError(err); ok {
		response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
		return
	}
	response.SendErrorResponse(w, http.StatusInternalServerError, fallbackCode, err, fallbackMsg)
}
