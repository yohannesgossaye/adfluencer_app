/**
 * One word list for every status the API returns, so each status reads the
 * same on every screen. Render with <StatusBadge> (components/common/status-badge).
 */
export type StatusTone = "success" | "warning" | "danger" | "info" | "neutral";
export type StatusInfo = { label: string; tone: StatusTone };

const fallback = (status: string): StatusInfo => ({
  label: status ? status.charAt(0).toUpperCase() + status.slice(1).replace(/_/g, " ") : "Unknown",
  tone: "neutral",
});
const lookup = (map: Record<string, StatusInfo>) => (status?: string | null): StatusInfo =>
  (status && map[status.toLowerCase()]) || fallback(status ?? "");

export const campaignStatus = lookup({
  draft: { label: "Draft", tone: "neutral" },
  pending: { label: "Draft", tone: "neutral" },
  active: { label: "Live", tone: "success" },
  paused: { label: "Paused by Adfluence", tone: "warning" },
  completed: { label: "Completed", tone: "info" },
  cancelled: { label: "Cancelled", tone: "neutral" },
});

/** `as` picks the wording for who is looking at the invite. */
export function inviteStatus(status?: string | null, as: "business" | "influencer" = "business"): StatusInfo {
  const map: Record<string, StatusInfo> = {
    pending: { label: as === "business" ? "Waiting for reply" : "Needs your answer", tone: "warning" },
    accepted: { label: "Accepted", tone: "success" },
    rejected: { label: as === "business" ? "Declined" : "You declined", tone: "neutral" },
    cancelled: { label: "Withdrawn", tone: "neutral" },
    expired: { label: "Expired", tone: "neutral" },
  };
  return (status && map[status.toLowerCase()]) || fallback(status ?? "");
}

/** A creator's account: profile review state. */
export const creatorStatus = lookup({
  pending_profile: { label: "Profile not submitted", tone: "neutral" },
  pending_approval: { label: "Waiting for review", tone: "warning" },
  pending: { label: "Waiting for review", tone: "warning" },
  active: { label: "Approved", tone: "success" },
  approved: { label: "Approved", tone: "success" },
  rejected: { label: "Changes requested", tone: "danger" },
  suspended: { label: "Suspended", tone: "danger" },
  blocked: { label: "Blocked", tone: "danger" },
  removed: { label: "Removed", tone: "danger" },
});

export const businessVerificationStatus = lookup({
  unsubmitted: { label: "Not submitted", tone: "neutral" },
  pending: { label: "Waiting for review", tone: "warning" },
  verified: { label: "Verified", tone: "success" },
  rejected: { label: "Not approved", tone: "danger" },
});

/** Any account in the admin Users list. */
export const accountStatus = lookup({
  active: { label: "Active", tone: "success" },
  pending: { label: "Pending", tone: "warning" },
  pending_profile: { label: "Profile not submitted", tone: "neutral" },
  pending_approval: { label: "Waiting for review", tone: "warning" },
  rejected: { label: "Changes requested", tone: "danger" },
  suspended: { label: "Suspended", tone: "danger" },
  blocked: { label: "Blocked", tone: "danger" },
  removed: { label: "Removed", tone: "danger" },
});

export const disputeStatus = lookup({
  open: { label: "Open", tone: "warning" },
  under_review: { label: "Under review", tone: "info" },
  escalated: { label: "Escalated", tone: "danger" },
  resolved: { label: "Resolved", tone: "success" },
  closed: { label: "Closed", tone: "neutral" },
});

export const ticketStatus = lookup({
  open: { label: "Open", tone: "warning" },
  in_review: { label: "In review", tone: "info" },
  escalated: { label: "Escalated", tone: "danger" },
  resolved: { label: "Resolved", tone: "success" },
  closed: { label: "Closed", tone: "neutral" },
});

/** Dispute severity and ticket priority share one scale. */
export const priorityLevel = lookup({
  low: { label: "Low", tone: "neutral" },
  normal: { label: "Normal", tone: "neutral" },
  medium: { label: "Medium", tone: "warning" },
  high: { label: "High", tone: "danger" },
  urgent: { label: "Urgent", tone: "danger" },
});

export const PLATFORM_LABEL: Record<string, string> = {
  tiktok: "TikTok",
  instagram: "Instagram",
  youtube: "YouTube",
  telegram: "Telegram",
  facebook: "Facebook",
};

export function platformLabel(platform?: string | null) {
  if (!platform) return "";
  return PLATFORM_LABEL[platform.toLowerCase()] ?? platform;
}
