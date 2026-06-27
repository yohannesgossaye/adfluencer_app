export type CreatorChatRole = "business" | "influencer" | "system";

export type CreatorChatMessage = {
  id: string;
  role: CreatorChatRole;
  body: string;
  sentAt: string;
  attachmentLabel?: string;
};

export type CreatorConversation = {
  id: string;
  businessName: string;
  brandInitial: string;
  campaignName: string;
  lastMessage: string;
  lastTimeIso: string;
  unreadCount: number;
  messages: CreatorChatMessage[];
};

const iso = (d: string, t: string) => new Date(`${d}T${t}:00+03:00`).toISOString();

export const initialCreatorConversations: CreatorConversation[] = [
  {
    id: "ic-1",
    businessName: "Betu Cosmetics",
    brandInitial: "B",
    campaignName: "Betu Cosmetics Launch",
    lastMessage: "Can you confirm delivery by Friday?",
    lastTimeIso: iso("2026-05-11", "09:15"),
    unreadCount: 1,
    messages: [
      {
        id: "ic1-1",
        role: "business",
        body: "Hi Biruk — excited to have you on the Betu launch.",
        sentAt: iso("2026-05-09", "10:00"),
      },
      {
        id: "ic1-2",
        role: "influencer",
        body: "Thanks! I can deliver the TikTok draft by Wednesday.",
        sentAt: iso("2026-05-09", "14:20"),
      },
      {
        id: "ic1-3",
        role: "business",
        body: "Can you confirm delivery by Friday?",
        sentAt: iso("2026-05-11", "09:15"),
      },
    ],
  },
  {
    id: "ic-2",
    businessName: "Habesha Coffee Roasters",
    brandInitial: "H",
    campaignName: "Coffee House Awareness",
    lastMessage: "The B-roll locations look perfect.",
    lastTimeIso: iso("2026-05-10", "16:00"),
    unreadCount: 0,
    messages: [
      {
        id: "ic2-1",
        role: "business",
        body: "Sharing the mood board for the café shoot.",
        sentAt: iso("2026-05-08", "11:00"),
        attachmentLabel: "moodboard.pdf",
      },
      {
        id: "ic2-2",
        role: "influencer",
        body: "Received — I’ll scout two spots in Bole.",
        sentAt: iso("2026-05-08", "15:30"),
      },
      {
        id: "ic2-3",
        role: "business",
        body: "The B-roll locations look perfect.",
        sentAt: iso("2026-05-10", "16:00"),
      },
    ],
  },
  {
    id: "ic-3",
    businessName: "Ride Ethiopia",
    brandInitial: "R",
    campaignName: "Ride App Creator Push",
    lastMessage: "Let’s add a CTA card at the end of the Short.",
    lastTimeIso: iso("2026-05-07", "12:40"),
    unreadCount: 0,
    messages: [
      {
        id: "ic3-1",
        role: "influencer",
        body: "First cut of the YouTube segment is ready for review.",
        sentAt: iso("2026-05-06", "18:00"),
      },
      {
        id: "ic3-2",
        role: "business",
        body: "Let’s add a CTA card at the end of the Short.",
        sentAt: iso("2026-05-07", "12:40"),
      },
    ],
  },
  {
    id: "ic-4",
    businessName: "Abyssinia Restaurant",
    brandInitial: "A",
    campaignName: "Abyssinia Restaurant Promo",
    lastMessage: "Weekend promo starts this Friday.",
    lastTimeIso: iso("2026-05-05", "08:00"),
    unreadCount: 0,
    messages: [
      {
        id: "ic4-1",
        role: "business",
        body: "Weekend promo starts this Friday.",
        sentAt: iso("2026-05-05", "08:00"),
      },
      {
        id: "ic4-2",
        role: "influencer",
        body: "On it — I’ll prioritize the menu hero shots.",
        sentAt: iso("2026-05-05", "10:10"),
      },
    ],
  },
];
