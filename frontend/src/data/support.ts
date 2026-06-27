import { initialAdminCampaigns } from "@/data/admin-campaigns-mock";

export type SupportRole = "business" | "influencer";

export type SupportPriority = "low" | "medium" | "high";

export type SupportStatus = "open" | "in_review" | "resolved" | "escalated";

export type SupportTicket = {
  id: string;
  senderName: string;
  role: SupportRole;
  email: string;
  subject: string;
  message: string;
  type: string;
  priority: SupportPriority;
  status: SupportStatus;
  date: string;
  relatedCampaign?: string;
  adminReply?: string;
  adminNotes?: string;
  source?: "seed" | "user";
};

export const SUPPORT_TICKETS_KEY = "adfluence_support_tickets";

const LEGACY_SUPPORT_KEY = "adfluence_admin_support";

export const businessIssueTypeOptions = [
  "Campaign issue",
  "Contract issue",
  "Payment question",
  "Creator communication",
  "Technical issue",
  "General feedback",
] as const;

export const influencerIssueTypeOptions = [
  "Profile approval",
  "Campaign invite",
  "Contract issue",
  "Payment/payout question",
  "Business communication",
  "Technical issue",
  "General feedback",
] as const;

export const relatedBusinessCampaignOptions = [
  "Not related to a specific campaign",
  ...initialAdminCampaigns.slice(0, 6).map((c) => c.title),
];

export const relatedInfluencerCampaignOptions = [
  "Not related to a specific invite",
  "Betu Cosmetics Launch — invite",
  "Ride App Creator Push — invite",
  "Summer Skincare Promo — invite",
  "Addis Food Week — invite",
];

const seedSupportTickets: SupportTicket[] = [
  {
    id: "tkt-b1",
    senderName: "Meron Gebre",
    role: "business",
    email: "meron@addisbrands.et",
    subject: "Creator stopped replying",
    message:
      "We had an active thread about deliverables and the creator has not responded in 6 days. What is the recommended next step?",
    type: "Creator communication",
    priority: "high",
    status: "open",
    date: "2026-05-10",
    relatedCampaign: "Betu Cosmetics Launch",
    source: "seed",
  },
  {
    id: "tkt-b2",
    senderName: "Helina Imports",
    role: "business",
    email: "ops@helinimports.et",
    subject: "Need help changing campaign budget",
    message:
      "We need to raise the budget mid-flight for our Ride App push. Is there a safe way to do this without breaking existing invites?",
    type: "Campaign issue",
    priority: "medium",
    status: "in_review",
    date: "2026-05-09",
    relatedCampaign: "Ride App Creator Push",
    adminNotes: "Checked escrow rules — follow up with finance mock.",
    source: "seed",
  },
  {
    id: "tkt-b3",
    senderName: "Yohannes",
    role: "business",
    email: "yohannes@gmail.com",
    subject: "Contract payment question",
    message:
      "A contract shows as funded but the milestone payment UI is greyed out. Can someone clarify the release timeline?",
    type: "Payment question",
    priority: "medium",
    status: "open",
    date: "2026-05-08",
    relatedCampaign: "Not related to a specific campaign",
    source: "seed",
  },
  {
    id: "tkt-inf1",
    senderName: "Biruk",
    role: "influencer",
    email: "biruk@gmail.com",
    subject: "Profile approval pending",
    message:
      "Submitted verification and portfolio links last week. Status still shows pending — is anything missing on my side?",
    type: "Profile approval",
    priority: "high",
    status: "open",
    date: "2026-05-07",
    relatedCampaign: "Not related to a specific invite",
    source: "seed",
  },
  {
    id: "tkt-inf2",
    senderName: "Kidus Alemayehu",
    role: "influencer",
    email: "kidus@gmail.com",
    subject: "Brand changed deliverables",
    message:
      "After signing, the business asked for two extra stories and a live session not in the original brief. How should I handle this?",
    type: "Campaign invite",
    priority: "high",
    status: "in_review",
    date: "2026-05-06",
    relatedCampaign: "Summer Skincare Promo — invite",
    source: "seed",
  },
  {
    id: "tkt-inf3",
    senderName: "Saron Tesfaye",
    role: "influencer",
    email: "saron.creates@gmail.com",
    subject: "Payout timeline question",
    message:
      "Completed all deliverables 10 days ago. When should I expect the payout to hit my wallet in the app?",
    type: "Payment/payout question",
    priority: "low",
    status: "resolved",
    date: "2026-04-20",
    relatedCampaign: "Betu Cosmetics Launch — invite",
    adminReply: "Payouts batch every Tuesday (demo). You should see it in the next cycle.",
    source: "seed",
  },
  {
    id: "tkt-b4",
    senderName: "Yohannes",
    role: "business",
    email: "yohannes@gmail.com",
    subject: "Invoice export not loading",
    message: "CSV export spins forever on the contracts page after filtering by active status.",
    type: "Technical issue",
    priority: "medium",
    status: "open",
    date: "2026-05-11",
    source: "seed",
  },
  {
    id: "tkt-b5",
    senderName: "Ride Ethiopia",
    role: "business",
    email: "partners@ride.et",
    subject: "Dispute with creator",
    message: "We want to escalate a quality issue from last month regarding missed posting windows.",
    type: "General feedback",
    priority: "high",
    status: "escalated",
    date: "2026-04-22",
    relatedCampaign: "Ride App Creator Push",
    source: "seed",
  },
  {
    id: "tkt-inf4",
    senderName: "Liya Bekele",
    role: "influencer",
    email: "liya.content@gmail.com",
    subject: "Payout schedule",
    message: "When is the next batch payout for completed contracts in March?",
    type: "Payment/payout question",
    priority: "medium",
    status: "resolved",
    date: "2026-05-01",
    adminReply: "March batch was released 2026-04-28 (mock).",
    source: "seed",
  },
];

export const initialSupportTickets = seedSupportTickets;

function isSupportStatus(v: unknown): v is SupportStatus {
  return v === "open" || v === "in_review" || v === "resolved" || v === "escalated";
}

function normalizeTicket(raw: unknown): SupportTicket | null {
  if (!raw || typeof raw !== "object") return null;
  const r = raw as Record<string, unknown>;
  if (typeof r.id !== "string") return null;
  const role: SupportRole = r.role === "influencer" ? "influencer" : "business";
  let status: SupportStatus = "open";
  if (isSupportStatus(r.status)) status = r.status;
  let priority: SupportPriority = "medium";
  if (r.priority === "low" || r.priority === "medium" || r.priority === "high") priority = r.priority;

  return {
    id: r.id,
    senderName: typeof r.senderName === "string" ? r.senderName : "Unknown",
    role,
    email: typeof r.email === "string" ? r.email : "",
    subject: typeof r.subject === "string" ? r.subject : "",
    message: typeof r.message === "string" ? r.message : "",
    type: typeof r.type === "string" ? r.type : "General feedback",
    priority,
    status,
    date: typeof r.date === "string" ? r.date : new Date().toISOString().slice(0, 10),
    relatedCampaign: typeof r.relatedCampaign === "string" ? r.relatedCampaign : undefined,
    adminReply: typeof r.adminReply === "string" ? r.adminReply : undefined,
    adminNotes: typeof r.adminNotes === "string" ? r.adminNotes : undefined,
    source: r.source === "user" || r.source === "seed" ? r.source : undefined,
  };
}

function tryParseTickets(raw: string | null): SupportTicket[] | null {
  if (!raw) return null;
  try {
    const parsed = JSON.parse(raw) as unknown;
    if (!Array.isArray(parsed)) return null;
    const rows = parsed.map(normalizeTicket).filter(Boolean) as SupportTicket[];
    return rows.length ? rows : null;
  } catch {
    return null;
  }
}

function migrateLegacyStorage(): SupportTicket[] | null {
  if (typeof window === "undefined") return null;
  const legacy = tryParseTickets(localStorage.getItem(LEGACY_SUPPORT_KEY));
  if (!legacy) return null;
  localStorage.setItem(SUPPORT_TICKETS_KEY, JSON.stringify(legacy));
  localStorage.removeItem(LEGACY_SUPPORT_KEY);
  return legacy;
}

export function loadSupportTickets(): SupportTicket[] {
  if (typeof window === "undefined") return [...seedSupportTickets];
  const fromNew = tryParseTickets(localStorage.getItem(SUPPORT_TICKETS_KEY));
  if (fromNew) return fromNew;
  const migrated = migrateLegacyStorage();
  if (migrated) return migrated;
  return [...seedSupportTickets];
}

export function saveSupportTickets(rows: SupportTicket[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(SUPPORT_TICKETS_KEY, JSON.stringify(rows));
}

export function appendSupportTicket(ticket: SupportTicket) {
  const rows = [...loadSupportTickets()];
  rows.push(ticket);
  saveSupportTickets(rows);
}

export function filterTicketsForPortalUser(
  tickets: SupportTicket[],
  role: SupportRole,
  userEmail: string | null,
): SupportTicket[] {
  if (!userEmail) return [];
  return tickets.filter((t) => t.role === role && t.email.toLowerCase() === userEmail.toLowerCase());
}
