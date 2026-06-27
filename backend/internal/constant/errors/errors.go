package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerErrror       = errors.New("internal server error")
	ErrBadRequest                 = errors.New("bad request")
	ErrNameIsRequired             = errors.New("name is required")
	ErrInvalidNameLength          = errors.New("invalid name length")
	ErrEmailIsRequired            = errors.New("email is required")
	ErrInvalidEmail               = errors.New("email must be a valid email address")
	ErrPhoneIsRequired            = errors.New("phone number is required")
	ErrInvalidPhone               = errors.New("phone number must be a valid Ethiopian phone number")
	ErrRoleIsRequired             = errors.New("role is required")
	ErrInvalidRole                = errors.New("role must be either business or influencer")
	ErrPasswordIsRequired         = errors.New("password is required")
	ErrInvalidPasswordLength      = errors.New("password must be between 8 and 20 characters")
	ErrConfirmPasswordRequired    = errors.New("confirm password is required")
	ErrPasswordMismatch           = errors.New("confirm password must match password")
	ErrEmailAlreadyExists         = errors.New("email already registered")
	ErrPhoneAlreadyExists         = errors.New("phone number already registered")
	ErrInvalidCredentials         = errors.New("incorrect email or password")
	ErrAccountBlocked             = errors.New("this account has been blocked")
	ErrAccountSuspended           = errors.New("account suspended")
	ErrAccountRemoved             = errors.New("account no longer available")
	ErrNotAdmin                   = errors.New("admin access required")
	ErrInvalidPortal              = errors.New("portal must be empty or admin")
	ErrUnauthorized               = errors.New("unauthorized")
	ErrRefreshTokenRequired       = errors.New("refresh token is required")
	ErrEmailAlreadyVerified       = errors.New("email is already verified")
	ErrTooManyVerifyEmailRequests = errors.New("too many verification emails requested")
	ErrPhoneAlreadyVerified       = errors.New("phone is already verified")
	ErrTooManyVerifyPhoneRequests = errors.New("too many verification sms requested")
	ErrNoPhoneOnAccount           = errors.New("no phone number on account")
	ErrProfileNotFound            = errors.New("influencer profile not found")
	ErrProfileIncomplete          = errors.New("profile incomplete")
	ErrInvalidStatusForSubmit     = errors.New("invalid status for submit")
	ErrProfileLocked              = errors.New("profile locked")
	ErrForbiddenRole              = errors.New("forbidden")
	ErrUnknownCategory            = errors.New("unknown category")
	ErrApplicationNotFound        = errors.New("application not found")
	ErrInvalidStatusForReview     = errors.New("application is not pending approval")
	ErrStorageNotConfigured       = errors.New("storage is not configured")
	ErrInvalidFile                = errors.New("invalid file")
	ErrPortfolioLimitReached      = errors.New("portfolio limit reached")
	ErrPortfolioItemNotFound      = errors.New("portfolio item not found")
	ErrLoginLocked                = errors.New("too many failed login attempts")
	ErrForgotPasswordRateLimit    = errors.New("too many password reset requests")
	ErrInvalidStatusForComplete   = errors.New("invalid status for complete")
	ErrAccountNotVerified         = errors.New("account not verified")
	ErrCampaignNotFound           = errors.New("campaign not found")
	ErrInvalidCampaignStatus      = errors.New("invalid campaign status")
	ErrAccountNotActive           = errors.New("account not active")
	ErrCreatorNotFound            = errors.New("creator not found")
	ErrInviteNotFound             = errors.New("invite not found")
	ErrInviteAlreadyExists        = errors.New("invite already exists")
	ErrInvalidInviteStatus        = errors.New("invalid invite status")
	ErrContractNotFound           = errors.New("contract not found")
	ErrInvalidContractStatus      = errors.New("invalid contract status")
	ErrInviteNotAccepted          = errors.New("invite not accepted")
	ErrContractAlreadyExists      = errors.New("contract already exists")
	ErrPaymentAlreadyExists       = errors.New("payment already exists")
	ErrPaymentRequired            = errors.New("payment required")
	ErrBusinessNotFound           = errors.New("business not found")
	ErrBusinessAlreadyVerified    = errors.New("business already verified")
	ErrBusinessNotVerified        = errors.New("business is not verified")
	ErrSupportTicketNotFound      = errors.New("support ticket not found")
	ErrCategoryNotFound           = errors.New("category not found")
	ErrCategoryAlreadyExists      = errors.New("category already exists")
	ErrUserNotFound               = errors.New("user not found")
	ErrCannotModifyAdminUser      = errors.New("cannot modify admin user")
	ErrInvalidUserStatusChange    = errors.New("invalid user status change")
	ErrDisputeNotFound            = errors.New("dispute not found")
	ErrInvalidDisputeStatus       = errors.New("invalid dispute status")
)

var ErrorMap = map[error]int{
	ErrInternalServerErrror:       http.StatusInternalServerError,
	ErrBadRequest:                 http.StatusBadRequest,
	ErrNameIsRequired:             http.StatusBadRequest,
	ErrInvalidNameLength:          http.StatusBadRequest,
	ErrEmailIsRequired:            http.StatusBadRequest,
	ErrInvalidEmail:               http.StatusBadRequest,
	ErrPhoneIsRequired:            http.StatusBadRequest,
	ErrInvalidPhone:               http.StatusBadRequest,
	ErrRoleIsRequired:             http.StatusBadRequest,
	ErrInvalidRole:                http.StatusBadRequest,
	ErrPasswordIsRequired:         http.StatusBadRequest,
	ErrInvalidPasswordLength:      http.StatusBadRequest,
	ErrConfirmPasswordRequired:    http.StatusBadRequest,
	ErrPasswordMismatch:           http.StatusBadRequest,
	ErrEmailAlreadyExists:         http.StatusConflict,
	ErrPhoneAlreadyExists:         http.StatusConflict,
	ErrInvalidCredentials:         http.StatusUnauthorized,
	ErrAccountBlocked:             http.StatusForbidden,
	ErrAccountSuspended:           http.StatusForbidden,
	ErrAccountRemoved:             http.StatusForbidden,
	ErrNotAdmin:                   http.StatusForbidden,
	ErrInvalidPortal:              http.StatusBadRequest,
	ErrUnauthorized:               http.StatusUnauthorized,
	ErrRefreshTokenRequired:       http.StatusBadRequest,
	ErrEmailAlreadyVerified:       http.StatusConflict,
	ErrTooManyVerifyEmailRequests: http.StatusTooManyRequests,
	ErrPhoneAlreadyVerified:       http.StatusConflict,
	ErrTooManyVerifyPhoneRequests: http.StatusTooManyRequests,
	ErrNoPhoneOnAccount:           http.StatusBadRequest,
	ErrProfileNotFound:            http.StatusNotFound,
	ErrProfileIncomplete:          http.StatusBadRequest,
	ErrInvalidStatusForSubmit:     http.StatusConflict,
	ErrProfileLocked:              http.StatusForbidden,
	ErrForbiddenRole:              http.StatusForbidden,
	ErrUnknownCategory:            http.StatusBadRequest,
	ErrApplicationNotFound:        http.StatusNotFound,
	ErrInvalidStatusForReview:     http.StatusConflict,
	ErrStorageNotConfigured:       http.StatusServiceUnavailable,
	ErrInvalidFile:                http.StatusBadRequest,
	ErrPortfolioLimitReached:      http.StatusConflict,
	ErrPortfolioItemNotFound:      http.StatusNotFound,
	ErrLoginLocked:                http.StatusTooManyRequests,
	ErrForgotPasswordRateLimit:    http.StatusTooManyRequests,
	ErrInvalidStatusForComplete:   http.StatusConflict,
	ErrAccountNotVerified:         http.StatusForbidden,
	ErrCampaignNotFound:           http.StatusNotFound,
	ErrInvalidCampaignStatus:      http.StatusConflict,
	ErrAccountNotActive:           http.StatusForbidden,
	ErrCreatorNotFound:            http.StatusNotFound,
	ErrInviteNotFound:             http.StatusNotFound,
	ErrInviteAlreadyExists:        http.StatusConflict,
	ErrInvalidInviteStatus:        http.StatusConflict,
	ErrContractNotFound:           http.StatusNotFound,
	ErrInvalidContractStatus:      http.StatusConflict,
	ErrInviteNotAccepted:          http.StatusConflict,
	ErrContractAlreadyExists:      http.StatusConflict,
	ErrPaymentAlreadyExists:       http.StatusConflict,
	ErrPaymentRequired:            http.StatusConflict,
	ErrBusinessNotFound:           http.StatusNotFound,
	ErrBusinessAlreadyVerified:    http.StatusConflict,
	ErrBusinessNotVerified:      http.StatusConflict,
	ErrSupportTicketNotFound:    http.StatusNotFound,
	ErrCategoryNotFound:         http.StatusNotFound,
	ErrCategoryAlreadyExists:    http.StatusConflict,
	ErrUserNotFound:             http.StatusNotFound,
	ErrCannotModifyAdminUser:    http.StatusForbidden,
	ErrInvalidUserStatusChange:  http.StatusConflict,
	ErrDisputeNotFound:          http.StatusNotFound,
	ErrInvalidDisputeStatus:     http.StatusBadRequest,
}
