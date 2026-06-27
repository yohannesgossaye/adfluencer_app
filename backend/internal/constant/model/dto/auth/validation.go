package auth

import (
	"regexp"

	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	emailRegex          = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	ethiopianPhoneRegex = regexp.MustCompile(`^(\+251|251|0)?9[0-9]{8}$`)
)

func (r SignUpRequest) ValidateSignUpRequest() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error(customErr.ErrNameIsRequired.Error()),
			validation.Length(2, 50).Error(customErr.ErrInvalidNameLength.Error()),
		),

		validation.Field(&r.Email,
			validation.Required.Error(customErr.ErrEmailIsRequired.Error()),
			validation.Match(emailRegex).Error(customErr.ErrInvalidEmail.Error()),
		),

		validation.Field(&r.Phone,
			validation.Required.Error(customErr.ErrPhoneIsRequired.Error()),
			validation.Match(ethiopianPhoneRegex).Error(customErr.ErrInvalidPhone.Error()),
		),

		validation.Field(&r.Role,
			validation.Required.Error(customErr.ErrRoleIsRequired.Error()),
			validation.In(string(model.BusinessRole), string(model.InfluencerRole)).Error(customErr.ErrInvalidRole.Error()),
		),

		validation.Field(&r.Password,
			validation.Required.Error(customErr.ErrPasswordIsRequired.Error()),
			validation.Length(8, 20).Error(customErr.ErrInvalidPasswordLength.Error()),
		),

		validation.Field(&r.ConfirmPassword,
			validation.Required.Error(customErr.ErrConfirmPasswordRequired.Error()),
			validation.By(func(value interface{}) error {
				confirmPassword, _ := value.(string)
				if confirmPassword != r.Password {
					return validation.NewError(
						"validation_password_mismatch",
						customErr.ErrPasswordMismatch.Error(),
					)
				}
				return nil
			}),
		),
	)
}

func (r LoginRequest) ValidateLoginRequest() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email,
			validation.Required.Error(customErr.ErrEmailIsRequired.Error()),
			validation.Match(emailRegex).Error(customErr.ErrInvalidEmail.Error()),
		),
		validation.Field(&r.Password,
			validation.Required.Error(customErr.ErrPasswordIsRequired.Error()),
		),
		validation.Field(&r.Portal,
			validation.In("", "admin").Error(customErr.ErrInvalidPortal.Error()),
		),
	)
}

func (r LogoutRequest) ValidateLogoutRequest() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RefreshToken,
			validation.Required.Error(customErr.ErrRefreshTokenRequired.Error()),
		),
	)
}

func (r ForgotPasswordRequest) ValidateForgotPasswordRequest() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email,
			validation.Required.Error(customErr.ErrEmailIsRequired.Error()),
			validation.Match(emailRegex).Error(customErr.ErrInvalidEmail.Error()),
		),
	)
}

func (r ResetPasswordRequest) ValidateResetPasswordRequest() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token,
			validation.Required.Error(customErr.ErrBadRequest.Error()),
		),
		validation.Field(&r.Password,
			validation.Required.Error(customErr.ErrPasswordIsRequired.Error()),
			validation.Length(8, 20).Error(customErr.ErrInvalidPasswordLength.Error()),
		),
		validation.Field(&r.ConfirmPassword,
			validation.Required.Error(customErr.ErrConfirmPasswordRequired.Error()),
			validation.By(func(value interface{}) error {
				confirmPassword, _ := value.(string)
				if confirmPassword != r.Password {
					return validation.NewError(
						"validation_password_mismatch",
						customErr.ErrPasswordMismatch.Error(),
					)
				}
				return nil
			}),
		),
	)
}

func (r VerifyEmailRequest) ValidateVerifyEmailRequest() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OTP,
			validation.Required.Error(customErr.ErrBadRequest.Error()),
			validation.Length(6, 6).Error(customErr.ErrBadRequest.Error()),
		),
	)
}

func (r VerifyPhoneRequest) ValidateVerifyPhoneRequest() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OTP,
			validation.Required.Error(customErr.ErrBadRequest.Error()),
			validation.Length(6, 6).Error(customErr.ErrBadRequest.Error()),
		),
	)
}
