import type { StatItem } from "@/types";

export const businessStats: StatItem[] = [
  { id: "b1", label: "Active Campaigns", value: "8", hint: "+2 this week" },
  { id: "b2", label: "Recommended Influencers", value: "24" },
  { id: "b3", label: "Pending Contracts", value: "5" },
  { id: "b4", label: "Recent Messages", value: "13", hint: "4 unread" },
];

export const influencerStats: StatItem[] = [
  { id: "i1", label: "Profile Completion", value: "82%" },
  { id: "i2", label: "Campaign Invites", value: "6", hint: "2 new today" },
  { id: "i3", label: "Active Collaborations", value: "3" },
  { id: "i4", label: "Earnings Placeholder", value: "ETB 48,000" },
];

export const adminStats: StatItem[] = [
  { id: "a1", label: "Pending Influencer Approvals", value: "17" },
  { id: "a2", label: "Total Campaigns", value: "142" },
  { id: "a3", label: "Total Users", value: "2,304" },
  { id: "a4", label: "Recent Activity", value: "39 actions" },
];
