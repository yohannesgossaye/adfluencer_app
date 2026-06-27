export type DisputeSeverity = "low" | "medium" | "high";
export type DisputeStatus = "open" | "under_review" | "resolved";
export type IssueType =
  | "Payment issue"
  | "Contract disagreement"
  | "Late delivery"
  | "Inappropriate message"
  | "Quality complaint";

export type AdminDispute = {
  id: string;
  title: string;
  business: string;
  influencer: string;
  influencerEmail: string;
  campaign: string;
  issueType: IssueType;
  severity: DisputeSeverity;
  status: DisputeStatus;
  preview: string;
  messages: { id: string; from: string; body: string }[];
  adminNote: string;
  resolutionNote: string;
};

export const initialAdminDisputes: AdminDispute[] = [
  {
    id: "dsp-1",
    title: "Payout delay — Betu campaign",
    business: "Betu Cosmetics",
    influencer: "Hana Tesfaye",
    influencerEmail: "hana.creates@email.com",
    campaign: "Betu Cosmetics Launch",
    issueType: "Payment issue",
    severity: "high",
    status: "open",
    preview: "Creator claims milestone payment is 10 days late.",
    messages: [
      { id: "m1", from: "Hana", body: "The second milestone was due May 1 — still no payout." },
      { id: "m2", from: "Betu", body: "Finance is processing; should clear this week." },
    ],
    adminNote: "",
    resolutionNote: "",
  },
  {
    id: "dsp-2",
    title: "Deliverable scope dispute",
    business: "Ride Ethiopia",
    influencer: "Biruk Alemu",
    influencerEmail: "biruk@gmail.com",
    campaign: "Ride App Creator Push",
    issueType: "Contract disagreement",
    severity: "medium",
    status: "under_review",
    preview: "Disagreement on number of Shorts included.",
    messages: [
      { id: "m1", from: "Biruk", body: "Contract says 2 Shorts; brand asked for a third unpaid." },
    ],
    adminNote: "Reviewing signed PDF v2.",
    resolutionNote: "",
  },
  {
    id: "dsp-3",
    title: "Late delivery warning",
    business: "Abyssinia Restaurant",
    influencer: "Saron Bekele",
    influencerEmail: "saron.creates@email.com",
    campaign: "Abyssinia Restaurant Promo",
    issueType: "Late delivery",
    severity: "low",
    status: "open",
    preview: "Business requested expedited delivery for holiday weekend.",
    messages: [{ id: "m1", from: "Abyssinia", body: "We need the Reel by Thursday for the promo." }],
    adminNote: "",
    resolutionNote: "",
  },
  {
    id: "dsp-4",
    title: "Resolved — quality complaint",
    business: "Habesha Coffee Roasters",
    influencer: "Meron Dawit",
    influencerEmail: "meron.creates@email.com",
    campaign: "Coffee House Awareness",
    issueType: "Quality complaint",
    severity: "medium",
    status: "resolved",
    preview: "Brand unhappy with first cut; reshoot agreed.",
    messages: [{ id: "m1", from: "Meron", body: "Delivered revised Reel set — brand approved." }],
    adminNote: "Mediator session May 4.",
    resolutionNote: "Resolved with revised deliverables and partial credit.",
  },
];

export const DISPUTES_KEY = "adfluence_admin_disputes";

export function loadDisputes(): AdminDispute[] {
  if (typeof window === "undefined") return initialAdminDisputes;
  try {
    const raw = localStorage.getItem(DISPUTES_KEY);
    if (!raw) return initialAdminDisputes;
    const parsed = JSON.parse(raw) as AdminDispute[];
    if (!Array.isArray(parsed) || !parsed.length) return initialAdminDisputes;
    return parsed.map((d) => ({
      ...d,
      influencerEmail: d.influencerEmail ?? "",
      adminNote: d.adminNote ?? "",
      resolutionNote: d.resolutionNote ?? "",
    }));
  } catch {
    return initialAdminDisputes;
  }
}

export function saveDisputes(rows: AdminDispute[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(DISPUTES_KEY, JSON.stringify(rows));
}
