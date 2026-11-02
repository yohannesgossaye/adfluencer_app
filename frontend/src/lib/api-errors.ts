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
