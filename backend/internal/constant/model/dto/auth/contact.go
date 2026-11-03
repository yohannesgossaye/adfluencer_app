package auth

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ContactRequest is a message from the public Contact us page.
type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"` // business, creator or other
	Subject string `json:"subject"`
	Message string `json:"message"`
	// Website is a hidden field people never fill in; bots usually do.
	Website string `json:"website"`
}

func (r ContactRequest) Normalize() ContactRequest {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Role = strings.ToLower(strings.TrimSpace(r.Role))
	r.Subject = strings.TrimSpace(r.Subject)
	r.Message = strings.TrimSpace(r.Message)
	return r
}

func (r ContactRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("please enter your name"), validation.RuneLength(2, 100).Error("name must be 2–100 characters")),
		validation.Field(&r.Email, validation.Required.Error("please enter your email"), validation.Match(emailRegex).Error("please enter a valid email address")),
		validation.Field(&r.Role, validation.In("business", "creator", "other").Error("please choose who you are")),
		validation.Field(&r.Subject, validation.RuneLength(0, 150).Error("subject must be at most 150 characters")),
		validation.Field(&r.Message, validation.Required.Error("please write a message"), validation.RuneLength(10, 3000).Error("message must be 10–3000 characters")),
	)
}
