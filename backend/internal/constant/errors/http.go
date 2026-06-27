package errors

import (
	"errors"
	"net/http"
)

func HTTPStatusFromError(err error) (int, bool) {
	for mappedErr, status := range ErrorMap {
		if errors.Is(err, mappedErr) {
			return status, true
		}
	}
	return http.StatusInternalServerError, false
}

// IsClientError returns true for expected request/business errors (4xx).
func IsClientError(err error) bool {
	status, ok := HTTPStatusFromError(err)
	return ok && status < http.StatusInternalServerError
}

func ClientErrorCode(err error) string {
	switch {
	case errors.Is(err, ErrEmailAlreadyExists):
		return "email_already_exists"
	case errors.Is(err, ErrPhoneAlreadyExists):
		return "phone_already_exists"
	case errors.Is(err, ErrInvalidEmail):
		return "invalid_email"
	case errors.Is(err, ErrInvalidPhone):
		return "invalid_phone"
	case errors.Is(err, ErrPasswordMismatch):
		return "password_mismatch"
	case errors.Is(err, ErrInvalidRole):
		return "invalid_role"
	case errors.Is(err, ErrInvalidCredentials):
		return "invalid_credentials"
	case errors.Is(err, ErrAccountBlocked):
		return "account_blocked"
	case errors.Is(err, ErrAccountSuspended):
		return "account_suspended"
	case errors.Is(err, ErrAccountRemoved):
		return "account_removed"
	case errors.Is(err, ErrNotAdmin):
		return "not_admin"
	case errors.Is(err, ErrInvalidPortal):
		return "invalid_portal"
	case errors.Is(err, ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, ErrRefreshTokenRequired):
		return "refresh_token_required"
	case errors.Is(err, ErrEmailAlreadyVerified):
		return "email_already_verified"
	case errors.Is(err, ErrTooManyVerifyEmailRequests):
		return "too_many_verify_email_requests"
	case errors.Is(err, ErrPhoneAlreadyVerified):
		return "phone_already_verified"
	case errors.Is(err, ErrTooManyVerifyPhoneRequests):
		return "too_many_verify_phone_requests"
	case errors.Is(err, ErrNoPhoneOnAccount):
		return "no_phone_on_account"
	case errors.Is(err, ErrProfileNotFound):
		return "profile_not_found"
	case errors.Is(err, ErrProfileIncomplete):
		return "profile_incomplete"
	case errors.Is(err, ErrInvalidStatusForSubmit):
		return "invalid_status_for_submit"
	case errors.Is(err, ErrProfileLocked):
		return "profile_locked"
	case errors.Is(err, ErrForbiddenRole):
		return "forbidden"
	case errors.Is(err, ErrUnknownCategory):
		return "invalid_niche"
	case errors.Is(err, ErrApplicationNotFound):
		return "application_not_found"
	case errors.Is(err, ErrInvalidStatusForReview):
		return "invalid_status_for_review"
	case errors.Is(err, ErrStorageNotConfigured):
		return "storage_not_configured"
	case errors.Is(err, ErrInvalidFile):
		return "invalid_file"
	case errors.Is(err, ErrPortfolioLimitReached):
		return "portfolio_limit_reached"
	case errors.Is(err, ErrPortfolioItemNotFound):
		return "portfolio_item_not_found"
	case errors.Is(err, ErrInvalidStatusForComplete):
		return "invalid_status_for_complete"
	case errors.Is(err, ErrAccountNotVerified):
		return "account_not_verified"
	case errors.Is(err, ErrCampaignNotFound):
		return "campaign_not_found"
	case errors.Is(err, ErrInvalidCampaignStatus):
		return "invalid_campaign_status"
	case errors.Is(err, ErrAccountNotActive):
		return "account_not_active"
	case errors.Is(err, ErrCreatorNotFound):
		return "creator_not_found"
	case errors.Is(err, ErrInviteNotFound):
		return "invite_not_found"
	case errors.Is(err, ErrInviteAlreadyExists):
		return "invite_already_exists"
	case errors.Is(err, ErrInvalidInviteStatus):
		return "invalid_invite_status"
	case errors.Is(err, ErrContractNotFound):
		return "contract_not_found"
	case errors.Is(err, ErrInvalidContractStatus):
		return "invalid_contract_status"
	case errors.Is(err, ErrInviteNotAccepted):
		return "invite_not_accepted"
	case errors.Is(err, ErrContractAlreadyExists):
		return "contract_already_exists"
	case errors.Is(err, ErrPaymentAlreadyExists):
		return "payment_already_exists"
	case errors.Is(err, ErrPaymentRequired):
		return "payment_required"
	case errors.Is(err, ErrBusinessNotFound):
		return "business_not_found"
	case errors.Is(err, ErrBusinessAlreadyVerified):
		return "business_already_verified"
	case errors.Is(err, ErrBusinessNotVerified):
		return "business_not_verified"
	case errors.Is(err, ErrSupportTicketNotFound):
		return "support_ticket_not_found"
	case errors.Is(err, ErrCategoryNotFound):
		return "category_not_found"
	case errors.Is(err, ErrCategoryAlreadyExists):
		return "category_already_exists"
	case errors.Is(err, ErrUserNotFound):
		return "user_not_found"
	case errors.Is(err, ErrCannotModifyAdminUser):
		return "cannot_modify_admin_user"
	case errors.Is(err, ErrInvalidUserStatusChange):
		return "invalid_user_status_change"
	case errors.Is(err, ErrDisputeNotFound):
		return "dispute_not_found"
	case errors.Is(err, ErrInvalidDisputeStatus):
		return "invalid_dispute_status"
	case errors.Is(err, ErrLoginLocked):
		return "login_locked"
	case errors.Is(err, ErrForgotPasswordRateLimit):
		return "too_many_forgot_password_requests"
	default:
		return "bad_request"
	}
}
