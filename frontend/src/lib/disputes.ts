/** Dispute vocabulary. Keep in step with backend dto/contract/dispute.go. */

export type DisputeRole = "business" | "influencer";
export type DisputeOutcome = "continue" | "complete" | "terminate";

export const DISPUTE_REASONS: Record<DisputeRole, Array<{ id: string; label: string }>> = {
  business: [
    { id: "not_delivered", label: "Work not delivered" },
    { id: "delivered_late", label: "Delivered late" },
    { id: "not_as_agreed", label: "Work not as agreed" },
    { id: "creator_unresponsive", label: "Creator not responding" },
    { id: "other", label: "Other" },
  ],
  influencer: [
    { id: "business_unresponsive", label: "Business not responding" },
    { id: "extra_work_requested", label: "Asked for extra work not in the contract" },
    { id: "approval_refused", label: "Approval unfairly refused" },
    { id: "other", label: "Other" },
  ],
};

export const DESCRIPTION_MIN = 20;
export const DESCRIPTION_MAX = 2000;
export const DECISION_NOTE_MIN = 10;

/** What each admin decision does, in the words both sides see. */
export const DISPUTE_OUTCOMES: Record<DisputeOutcome, { label: string; effect: string; tone: "default" | "danger" }> = {
  continue: {
    label: "Continue the contract",
    effect: "The contract goes back to in progress. The payment stays held until the business approves the work.",
    tone: "default",
  },
  complete: {
    label: "Complete and pay the creator",
    effect: "The work is accepted, the contract is completed and the creator's payout is released.",
    tone: "default",
  },
  terminate: {
    label: "End the contract",
    effect: "The contract is cancelled and Adfluence refunds the business.",
    tone: "danger",
  },
};

export function outcomeLabel(outcome?: string) {
  return outcome && outcome in DISPUTE_OUTCOMES ? DISPUTE_OUTCOMES[outcome as DisputeOutcome].label : "Resolved";
}
