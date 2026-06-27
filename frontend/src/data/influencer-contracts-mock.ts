export type InfluencerContractStatus =
  | "pending_review"
  | "pending"
  | "accepted"
  | "active"
  | "rejected"
  | "completed";

export type InfluencerContractRecord = {
  id: string;
  brand: string;
  campaign: string;
  deliverables: string;
  amountEtb: number;
  timeline: string;
  status: InfluencerContractStatus;
  terms: string;
  deadline: string;
  notes: string;
  checklist: { id: string; label: string; done: boolean }[];
};

export const influencerContractStats = {
  pending: 2,
  active: 2,
  completed: 8,
  totalEarningsEtb: 96000,
} as const;

export const initialInfluencerContracts: InfluencerContractRecord[] = [
  {
    id: "ict-1",
    brand: "Betu Cosmetics",
    campaign: "Betu Cosmetics Launch",
    deliverables: "2 TikToks, 3 stories, 1 Reel",
    amountEtb: 280000,
    timeline: "May 12 – Jun 2, 2026",
    status: "active",
    terms: "90-day organic usage; paid ads optional add-on.",
    deadline: "Jun 2, 2026",
    notes: "Primary contact: campaigns@betu.et",
    checklist: [
      { id: "c1", label: "TikTok 1 delivered", done: true },
      { id: "c2", label: "Stories live", done: false },
    ],
  },
  {
    id: "ict-2",
    brand: "Ride Ethiopia",
    campaign: "Ride App Creator Push",
    deliverables: "1 YouTube, 2 Shorts",
    amountEtb: 340000,
    timeline: "May 15 – Jun 10, 2026",
    status: "pending_review",
    terms: "Creator retains moral rights; brand gets 12-month license.",
    deadline: "Jun 10, 2026",
    notes: "Awaiting your signature.",
    checklist: [
      { id: "c1", label: "Contract signed", done: false },
      { id: "c2", label: "Deposit received", done: false },
    ],
  },
  {
    id: "ict-3",
    brand: "Habesha Coffee Roasters",
    campaign: "Coffee House Awareness",
    deliverables: "2 Reels, 4 stories",
    amountEtb: 190000,
    timeline: "May 8 – May 28, 2026",
    status: "pending",
    terms: "Standard Adfluence creator agreement.",
    deadline: "May 28, 2026",
    notes: "Please review deliverable dates.",
    checklist: [
      { id: "c1", label: "Kickoff call", done: true },
      { id: "c2", label: "First draft", done: false },
    ],
  },
  {
    id: "ict-4",
    brand: "Abyssinia Restaurant",
    campaign: "Abyssinia Restaurant Promo",
    deliverables: "3 TikToks, 1 Reel",
    amountEtb: 150000,
    timeline: "Apr 20 – May 10, 2026",
    status: "completed",
    terms: "Campaign closed — all usage rights fulfilled.",
    deadline: "May 10, 2026",
    notes: "Final payout processed.",
    checklist: [
      { id: "c1", label: "All deliverables", done: true },
      { id: "c2", label: "Final report", done: true },
    ],
  },
];
