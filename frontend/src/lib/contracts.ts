/** Shared contract vocabulary for the business and influencer portals. */
import type { StatusInfo, StatusTone } from "@/lib/status";

export const PLATFORM_FEE_PERCENT = 5; // charged to each side

export const DELIVERABLE_TYPES = [
  { id: "tiktok_video", label: "TikTok video" },
  { id: "instagram_reel", label: "Instagram reel" },
  { id: "instagram_story", label: "Instagram story" },
  { id: "youtube_short", label: "YouTube short" },
  { id: "product_review", label: "Product review" },
  { id: "post", label: "Post" },
  { id: "other", label: "Other" },
] as const;

export function deliverableLabel(type: string) {
  return DELIVERABLE_TYPES.find((d) => d.id === type)?.label ?? type.replace(/_/g, " ");
}

type Tone = "amber" | "sky" | "emerald" | "rose" | "slate";

const STATUS: Record<string, { label: string; business: string; influencer: string; tone: Tone }> = {
  pending_influencer: {
    label: "Waiting for creator",
    business: "Sent to the creator. They need to accept it first.",
    influencer: "Review the terms and accept or decline this contract.",
    tone: "amber",
  },
  accepted: {
    label: "Awaiting payment",
    business: "The creator accepted. Adfluence confirms your payment, then the creator starts.",
    influencer: "You accepted. Please wait: the work starts once Adfluence confirms the business's payment.",
    tone: "sky",
  },
  active: {
    label: "In progress",
    business: "Escrow is funded. Approve the work when the creator has delivered to release their payout.",
    influencer: "Escrow is funded. Deliver the work; the business approves it to release your payout.",
    tone: "sky",
  },
  completed: {
    label: "Completed",
    business: "Work approved and payout released to the creator.",
    influencer: "Work approved. Your payout has been released.",
    tone: "emerald",
  },
  cancelled: { label: "Cancelled", business: "This contract was cancelled.", influencer: "This contract was cancelled by the business.", tone: "slate" },
  rejected: { label: "Declined", business: "The creator declined this contract.", influencer: "You declined this contract.", tone: "slate" },
  disputed: { label: "In dispute", business: "Our team is reviewing a dispute on this contract.", influencer: "Our team is reviewing a dispute on this contract.", tone: "rose" },
};


// Badge wording from the creator's side where it differs from the business's.
const INFLUENCER_LABEL: Record<string, string> = {
  pending_influencer: "Needs your answer",
  accepted: "Waiting for payment",
  rejected: "You declined",
};

const BADGE_TONE: Record<Tone, StatusTone> = { amber: "warning", sky: "info", emerald: "success", rose: "danger", slate: "neutral" };

export function contractStatus(status: string, role: "business" | "influencer") {
  const s = STATUS[status] ?? { label: status.replace(/_/g, " "), business: "", influencer: "", tone: "slate" as Tone };
  const label = role === "influencer" ? (INFLUENCER_LABEL[status] ?? s.label) : s.label;
  /** `badge` is for <StatusBadge status={…} />. */
  return { label, hint: role === "business" ? s.business : s.influencer, badge: { label, tone: BADGE_TONE[s.tone] } as StatusInfo };
}

const round2 = (v: number) => Math.round(v * 100) / 100;

/** Same arithmetic as the backend (helper/fees.go): 5% from each side. */
export function previewFees(amount: number) {
  const fee = round2((amount * PLATFORM_FEE_PERCENT) / 100);
  return {
    amount,
    businessFee: fee,
    totalBusinessPays: round2(amount + fee),
    influencerFee: fee,
    influencerReceives: round2(amount - fee),
  };
}

/** Money is always written "ETB 12,500" across the app. */
export function formatEtb(value?: number) {
  return `ETB ${new Intl.NumberFormat("en-US", { maximumFractionDigits: 2 }).format(value ?? 0)}`;
}

/** "ETB 5,000 – 10,000"; open-ended when only one side is known. */
export function formatEtbRange(min?: number | null, max?: number | null) {
  const n = (v: number) => new Intl.NumberFormat("en-US", { maximumFractionDigits: 0 }).format(v);
  if (min && max) return `ETB ${n(min)} – ${n(max)}`;
  if (max) return `Up to ETB ${n(max)}`;
  if (min) return `From ETB ${n(min)}`;
  return "Price not set";
}

export function formatDate(iso?: string) {
  if (!iso) return "—";
  try {
    return new Intl.DateTimeFormat("en-GB", { day: "numeric", month: "short", year: "numeric" }).format(new Date(iso));
  } catch {
    return iso;
  }
}
