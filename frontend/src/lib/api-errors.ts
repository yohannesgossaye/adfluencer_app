import axios from "axios";

/** Plain-language text for error types the API returns (constant/errors/http.go, chat/errors.go). */
const BUSINESS_MESSAGES: Record<string, string> = {
  account_not_verified: "Your business needs to be verified by Adfluence first. Check your business profile.",
  business_not_verified: "Your business needs to be verified by Adfluence first. Check your business profile.",
  invalid_campaign_status: "This campaign isn't live. Publish it before inviting creators.",
  invite_already_exists: "You've already invited this creator to this campaign.",
  recipient_unavailable: "This creator can't receive messages right now.",
  chat_not_allowed: "Messaging isn't available for this account.",
  account_not_active: "Your account isn't active. Contact support.",
};

export function apiErrorType(err: unknown): string | undefined {
  if (!axios.isAxiosError(err)) return undefined;
  return (err.response?.data as { error?: { type?: string } } | undefined)?.error?.type;
}

/** A message a business user can act on, falling back to the server text, then `fallback`. */
export function businessErrorMessage(err: unknown, fallback: string): string {
  if (!axios.isAxiosError(err)) return fallback;
  if (!err.response) return "You appear to be offline. Please try again.";
  const type = apiErrorType(err);
  if (type && BUSINESS_MESSAGES[type]) return BUSINESS_MESSAGES[type];
  return (err.response.data as { message?: string } | undefined)?.message || fallback;
}

/** Plain-language text for errors an admin action can hit. */
const ADMIN_MESSAGES: Record<string, string> = {
  cannot_modify_admin_user: "Admin accounts can't be blocked or removed here.",
  invalid_user_status_change: "This account is already in that state. Refresh the list.",
  user_not_found: "This account no longer exists.",
  dispute_not_found: "This dispute no longer exists.",
  invalid_dispute_status: "This dispute has already moved on. Refresh the list.",
  campaign_not_found: "This campaign no longer exists.",
  invalid_campaign_status: "This campaign has already changed status. Refresh the list.",
  category_not_found: "This category no longer exists.",
  category_already_exists: "A category with this name already exists.",
  support_ticket_not_found: "This ticket no longer exists.",
  business_not_pending_review: "This business isn't waiting for review any more.",
  business_already_verified: "This business is already verified.",
  validation_error: "Some details are missing or invalid.",
  forbidden: "You don't have permission to do that.",
};

/** A message for an admin action, falling back to the server text, then `fallback`. */
export function adminErrorMessage(err: unknown, fallback: string): string {
  if (!axios.isAxiosError(err)) return fallback;
  if (!err.response) return "You appear to be offline. Please try again.";
  const type = apiErrorType(err);
  if (type && ADMIN_MESSAGES[type]) return ADMIN_MESSAGES[type];
  return (err.response.data as { message?: string } | undefined)?.message || fallback;
}

/** Plain-language text for login, signup and password errors. */
const AUTH_MESSAGES: Record<string, string> = {
  invalid_credentials: "That email and password don't match. Check them and try again.",
  login_locked: "Too many failed attempts. Wait a few minutes, or reset your password.",
  account_blocked: "This account is blocked. Contact Adfluence support.",
  account_suspended: "This account is suspended. Contact Adfluence support.",
  account_removed: "This account has been closed. Contact Adfluence support.",
  not_admin: "This login is for Adfluence staff. Business and creator accounts log in on the main login page.",
  email_already_exists: "An account with this email already exists. Log in instead.",
  phone_already_exists: "An account with this phone number already exists.",
  invalid_email: "Enter a valid email address.",
  invalid_phone: "Enter a valid phone number, e.g. +251 9XX XXX XXX.",
  password_mismatch: "The passwords don't match.",
  too_many_otp_attempts: "Too many wrong codes. Request a new code and try again.",
  too_many_forgot_password_requests: "Too many reset requests. Wait a few minutes and try again.",
  too_many_verify_email_requests: "Too many code requests. Wait a few minutes and try again.",
  validation_error: "Some details are missing or invalid.",
};

/** A message for the auth screens; never shows raw technical errors such as "Network Error". */
export function authErrorMessage(err: unknown, fallback: string): string {
  if (!axios.isAxiosError(err)) return fallback;
  if (!err.response) return "You appear to be offline. Check your connection and try again.";
  if (err.response.status >= 500) return "Adfluence is having trouble right now. Please try again in a moment.";
  const type = apiErrorType(err);
  if (type && AUTH_MESSAGES[type]) return AUTH_MESSAGES[type];
  const data = err.response.data as { message?: string; error?: { details?: Array<{ description?: string }> } } | undefined;
  return data?.error?.details?.[0]?.description || data?.message || fallback;
}
