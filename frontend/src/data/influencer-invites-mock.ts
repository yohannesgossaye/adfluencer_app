export type InfluencerInviteStatus = "pending" | "accepted" | "rejected";

export type InfluencerInvite = {
  id: string;
  brand: string;
  campaign: string;
  goal: string;
  budget: string;
  deadline: string;
  status: InfluencerInviteStatus;
  messagePreview: string;
  deliverables: string;
  timeline: string;
};

export const INVITES_STORAGE_KEY = "adfluence_influencer_invites";

export const initialInfluencerInvites: InfluencerInvite[] = [
  {
    id: "inv-1",
    brand: "Betu Cosmetics",
    campaign: "Betu Cosmetics Launch",
    goal: "Product launch",
    budget: "ETB 280,000",
    deadline: "Jun 2, 2026",
    status: "pending",
    messagePreview: "We’d love you to lead the hero TikTok for launch week.",
    deliverables: "2 TikToks, 3 stories, 1 Reel",
    timeline: "May 12 – Jun 2, 2026",
  },
  {
    id: "inv-2",
    brand: "Habesha Coffee Roasters",
    campaign: "Coffee House Awareness",
    goal: "Brand awareness",
    budget: "ETB 190,000",
    deadline: "May 28, 2026",
    status: "pending",
    messagePreview: "Café lifestyle content with a focus on Ethiopian coffee culture.",
    deliverables: "2 Reels, 4 stories",
    timeline: "May 8 – May 28, 2026",
  },
  {
    id: "inv-3",
    brand: "Ride Ethiopia",
    campaign: "Ride App Creator Push",
    goal: "App downloads",
    budget: "ETB 340,000",
    deadline: "Jun 10, 2026",
    status: "pending",
    messagePreview: "Creator-first walkthrough + Shorts cross-post.",
    deliverables: "1 YouTube feature, 2 Shorts",
    timeline: "May 15 – Jun 10, 2026",
  },
  {
    id: "inv-4",
    brand: "Abyssinia Restaurant",
    campaign: "Abyssinia Restaurant Promo",
    goal: "Event promotion",
    budget: "ETB 150,000",
    deadline: "May 18, 2026",
    status: "pending",
    messagePreview: "Weekend traffic push with menu highlights.",
    deliverables: "3 TikToks, 1 Reel",
    timeline: "Apr 20 – May 18, 2026",
  },
  {
    id: "inv-5",
    brand: "Sheba Tech Hub",
    campaign: "Campus Creator Challenge",
    goal: "Brand awareness",
    budget: "ETB 95,000",
    deadline: "Apr 30, 2026",
    status: "accepted",
    messagePreview: "Accepted — let’s align on filming dates.",
    deliverables: "1 longform, 2 Shorts",
    timeline: "Apr 1 – Apr 30, 2026",
  },
  {
    id: "inv-6",
    brand: "Legacy Bank",
    campaign: "Youth Savings Push",
    goal: "Sales conversion",
    budget: "ETB 120,000",
    deadline: "Mar 15, 2026",
    status: "rejected",
    messagePreview: "Declined due to schedule overlap.",
    deliverables: "Stories + Reel",
    timeline: "Feb 1 – Mar 15, 2026",
  },
];

export function loadInfluencerInvites(): InfluencerInvite[] {
  if (typeof window === "undefined") return initialInfluencerInvites;
  try {
    const raw = localStorage.getItem(INVITES_STORAGE_KEY);
    if (!raw) return initialInfluencerInvites;
    const parsed = JSON.parse(raw) as InfluencerInvite[];
    return Array.isArray(parsed) && parsed.length ? parsed : initialInfluencerInvites;
  } catch {
    return initialInfluencerInvites;
  }
}

export function saveInfluencerInvites(invites: InfluencerInvite[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(INVITES_STORAGE_KEY, JSON.stringify(invites));
}
