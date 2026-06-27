package response

import (
	"encoding/json"
	"net/http"
)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	detailMessage := message
	if err != nil {
		detailMessage = err.Error()
	}

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
