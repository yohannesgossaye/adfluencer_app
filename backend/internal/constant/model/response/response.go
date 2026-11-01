package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

// OnHiddenError, when set, receives errors whose details were withheld from
// the client (5xx responses and database errors) so they can still be logged.
var OnHiddenError func(status int, code string, err error)

func SendSuccessResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string, data, meta any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func SendErrorResponse(w http.ResponseWriter, status int, types string, err error, message string) {
	detailMessage := message
	if err != nil {
		detailMessage = err.Error()
	}

	if hideErrorDetails(status, err) {
		if OnHiddenError != nil && err != nil {
			OnHiddenError(status, types, err)
		}
		message = publicMessage(status)
		detailMessage = message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
		Error: &ErrorResponse{
			Type:    types,
			Message: message,
			Detail: []FieldError{
				{
					Title:       types,
					Description: detailMessage,
				},
			},
		},
	})
}

// hideErrorDetails reports whether the raw error text must not reach the
// client: server errors and anything coming from the database driver.
func hideErrorDetails(status int, err error) bool {
	if status >= http.StatusInternalServerError {
		return true
	}
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr)
}

func publicMessage(status int) string {
	if status >= http.StatusInternalServerError {
		return "internal server error"
	}
	return "request could not be processed"
}
