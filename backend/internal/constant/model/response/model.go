package response

type Response struct {
	Status  int            `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    any            `json:"data,omitempty"`
	Meta    any            `json:"meta,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Type    string       `json:"type,omitempty"`
	Message string       `json:"message,omitempty"`
	Detail  []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}