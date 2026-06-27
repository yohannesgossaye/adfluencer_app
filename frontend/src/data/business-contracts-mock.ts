export type BusinessContractRecordStatus =
  | "draft"
  | "pending"
  | "accepted"
  | "active"
  | "completed";

export type BusinessContractRecord = {
  id: string;
  title: string;
  influencerName: string;
  influencerImage: string;
  influencerId: string;
  campaign: string;
  deliverables: string;
  amountEtb: number;
  timeline: string;
  status: BusinessContractRecordStatus;
  brief: string;
  deliverableChecklist: { id: string; label: string; done: boolean }[];
  notes: string;
};

export const businessContractStats = {
  pendingReview: 3,
  active: 4,
  completed: 12,
  totalValueEtb: 185000,
} as const;

export const initialBusinessContracts: BusinessContractRecord[] = [
  {
    id: "bctr-1",
    title: "Betu Cosmetics x Hana Tesfaye",
    influencerName: "Hana Tesfaye",
    influencerImage: "/images/influencers/influencer-1.jpg",
    influencerId: "inf-1",
    campaign: "Betu Cosmetics Launch",
    deliverables: "2 TikToks, 3 IG stories, 1 Reel",
    amountEtb: 280000,
    timeline: "May 12 – Jun 2, 2026",
    status: "active",
    brief: "Product-led beauty campaign with authentic demos and before/after moments.",
    deliverableChecklist: [
      { id: "d1", label: "TikTok hero video", done: true },
      { id: "d2", label: "IG Reel", done: false },
      { id: "d3", label: "Story sequence", done: false },
    ],
    notes: "Usage rights: 90 days organic + whitelisting optional add-on.",
  },
  {
    id: "bctr-2",
    title: "Ride App Push x Biruk Alemu",
    influencerName: "Biruk Alemu",
    influencerImage: "/images/influencers/influencer-2.jpg",
    influencerId: "inf-2",
    campaign: "Ride App Creator Push",
    deliverables: "1 YouTube feature, 2 Shorts",
    amountEtb: 340000,
    timeline: "May 15 – Jun 10, 2026",
    status: "pending",
    brief: "Creator-first app walkthrough with CTA to download and first-ride promo.",
    deliverableChecklist: [
      { id: "d1", label: "YouTube longform", done: false },
      { id: "d2", label: "Shorts cut", done: false },
    ],
    notes: "Pending legal review on usage terms.",
  },
  {
    id: "bctr-3",
    title: "Coffee Awareness x Meron Dawit",
    influencerName: "Meron Dawit",
    influencerImage: "/images/influencers/influencer-3.jpg",
    influencerId: "inf-3",
    campaign: "Coffee House Awareness",
    deliverables: "2 Reels, 4 stories",
    amountEtb: 190000,
    timeline: "May 8 – May 28, 2026",
    status: "accepted",
    brief: "Café lifestyle storytelling with emphasis on Ethiopian coffee culture.",
    deliverableChecklist: [
      { id: "d1", label: "Reel 1", done: true },
      { id: "d2", label: "Reel 2", done: false },
    ],
    notes: "Accepted — awaiting signed PDF from creator.",
  },
  {
    id: "bctr-4",
    title: "Abyssinia Promo x Saron Bekele",
    influencerName: "Saron Bekele",
    influencerImage: "/images/influencers/influencer-4.jpg",
    influencerId: "inf-4",
    campaign: "Abyssinia Restaurant Promo",
    deliverables: "3 TikToks, 1 Reel",
    amountEtb: 150000,
    timeline: "Apr 20 – May 10, 2026",
    status: "completed",
    brief: "Weekend traffic push with menu highlights and ambiance.",
    deliverableChecklist: [
      { id: "d1", label: "TikTok series", done: true },
      { id: "d2", label: "Reel recap", done: true },
    ],
    notes: "Campaign closed — performance report shared.",
  },
  {
    id: "bctr-5",
    title: "Tech Week x Abel Girma",
    influencerName: "Abel Girma",
    influencerImage: "/images/influencers/influencer-5.jpg",
    influencerId: "inf-5",
    campaign: "Sheba Tech Week",
    deliverables: "1 review video, 2 Shorts",
    amountEtb: 120000,
    timeline: "Jun 1 – Jun 20, 2026",
    status: "draft",
    brief: "Hands-on product review with honest pros/cons.",
    deliverableChecklist: [
      { id: "d1", label: "Review video", done: false },
      { id: "d2", label: "Shorts", done: false },
    ],
    notes: "Draft — not sent to influencer yet.",
  },
];
