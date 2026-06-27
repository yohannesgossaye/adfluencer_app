export type AdminActivityItem = {
  id: string;
  type: "signup" | "campaign" | "contract" | "report";
  label: string;
  detail: string;
  time: string;
};

export const adminActivityFeed: AdminActivityItem[] = [
  {
    id: "act-1",
    type: "signup",
    label: "New signup",
    detail: "Betu Cosmetics — business verified",
    time: "2h ago",
  },
  {
    id: "act-2",
    type: "campaign",
    label: "Campaign created",
    detail: "Ride App Creator Push — Ride Ethiopia",
    time: "5h ago",
  },
  {
    id: "act-3",
    type: "contract",
    label: "Contract accepted",
    detail: "Hana Tesfaye × Betu Cosmetics Launch",
    time: "Yesterday",
  },
  {
    id: "act-4",
    type: "report",
    label: "Reported conversation",
    detail: "Flagged thread — payment dispute (review)",
    time: "Yesterday",
  },
];

export const adminAlerts = [
  {
    id: "al-1",
    title: "Reported chat",
    detail: "User reported inappropriate language in campaign thread.",
    severity: "high" as const,
  },
  {
    id: "al-2",
    title: "Delayed contract",
    detail: "Abyssinia Restaurant promo — milestone 2 overdue by 3 days.",
    severity: "medium" as const,
  },
  {
    id: "al-3",
    title: "Suspicious account",
    detail: "Multiple failed logins from new business account.",
    severity: "medium" as const,
  },
  {
    id: "al-4",
    title: "Payment issue",
    detail: "Payout batch #442 — manual verification suggested.",
    severity: "low" as const,
  },
];
