package httputil

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorMapping struct {
	Err    error
	Status int
	Code   string
}

// Postgres error codes that map to client errors rather than 500s.
const (
	pgInvalidTextRepresentation = "22P02" // e.g. malformed UUID in a path parameter
	pgUniqueViolation           = "23505"
)

func WriteServiceError(w http.ResponseWriter, err error, mappings ...ErrorMapping) {
	WriteServiceErrorOr(w, err, "internal_server_error", "internal server error", mappings...)
}

func WriteServiceErrorOr(w http.ResponseWriter, err error, fallbackCode, fallbackMsg string, mappings ...ErrorMapping) {
	var wait *customErr.WaitError
	if errors.As(err, &wait) {
		w.Header().Set("Retry-After", strconv.Itoa(wait.RetryAfterSeconds()))
	}
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
	if status, code, msg, ok := databaseClientError(err); ok {
		response.SendErrorResponse(w, status, code, nil, msg)
		return
	}
	response.SendErrorResponse(w, http.StatusInternalServerError, fallbackCode, err, fallbackMsg)
}

// databaseClientError translates database errors caused by bad client input
// into 4xx responses with a safe message.
func databaseClientError(err error) (int, string, string, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return 0, "", "", false
	}
	switch pgErr.Code {
	case pgInvalidTextRepresentation:
		return http.StatusBadRequest, "invalid_id", "invalid identifier", true
	case pgUniqueViolation:
		return http.StatusConflict, "conflict", "resource already exists", true
	}
	return 0, "", "", false
}
