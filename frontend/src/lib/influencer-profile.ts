/** Fields the server requires before a creator can submit their profile (helper/completeness.go). */
export const REQUIRED_FOR_SUBMIT: Record<string, string> = {
  creator_name: "Creator name",
  bio: "Bio of at least 20 characters",
  niche: "At least one niche",
  platform: "At least one platform",
  social_link: "A profile link for one of your platforms",
  followers: "Follower count",
  engagement_rate: "Engagement rate",
  price_range: "Minimum and maximum price",
  city: "City",
};

/** Optional items that raise the profile score but are not needed to submit. */
export const OPTIONAL_PROFILE_ITEMS: Record<string, string> = {
  languages: "Languages",
  profile_image: "Profile photo",
  portfolio_sample: "A portfolio sample",
};

export type CreatorStatus = "pending_profile" | "pending_approval" | "active" | "rejected" | "suspended" | "blocked" | "removed";

/** What a creator should read about their own account status. */
export const CREATOR_STATUS: Record<string, { label: string; tone: "neutral" | "amber" | "green" | "rose" }> = {
  pending_profile: { label: "Not submitted yet", tone: "neutral" },
  pending_approval: { label: "Under review", tone: "amber" },
  active: { label: "Approved", tone: "green" },
  rejected: { label: "Changes requested", tone: "rose" },
  suspended: { label: "Suspended", tone: "rose" },
  blocked: { label: "Blocked", tone: "rose" },
  removed: { label: "Removed", tone: "rose" },
};
