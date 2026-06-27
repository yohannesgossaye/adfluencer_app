package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r RejectApplicationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reason,
			validation.Required.Error("reason is required"),
			validation.Length(5, 1000).Error("reason must be between 5 and 1000 characters"),
		),
	)
}

func (q ListApplicationsQuery) Normalize() ListApplicationsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 20
	}
	if q.Limit > 50 {
		q.Limit = 50
	}
	q.Status = strings.TrimSpace(strings.ToLower(q.Status))
	q.Search = strings.TrimSpace(q.Search)
	return q
}

func (r UpdateAdminNotesRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AdminNote,
			validation.Required.Error(customErr.ErrBadRequest.Error()),
			validation.Length(1, 2000).Error("admin note must be between 1 and 2000 characters"),
		),
	)
}
