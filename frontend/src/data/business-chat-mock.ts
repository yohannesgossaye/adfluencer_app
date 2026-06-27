export type ChatMessageRole = "business" | "influencer" | "system";

export type BusinessChatMessage = {
  id: string;
  role: ChatMessageRole;
  body: string;
  sentAt: string;
  attachmentLabel?: string;
};

export type BusinessConversation = {
  id: string;
  influencerId: string;
  influencerName: string;
  imageUrl: string;
  campaignName: string;
  niche: string;
  platform: string;
  lastMessage: string;
  lastTimeIso: string;
  unreadCount: number;
  online: boolean;
  replyHint: string;
  messages: BusinessChatMessage[];
};

const iso = (d: string, t: string) => new Date(`${d}T${t}:00+03:00`).toISOString();

export const initialBusinessConversations: BusinessConversation[] = [
  {
    id: "conv-1",
    influencerId: "inf-1",
    influencerName: "Hana Tesfaye",
    imageUrl: "/images/influencers/influencer-1.jpg",
    campaignName: "Betu Cosmetics Launch",
    niche: "Beauty",
    platform: "TikTok",
    lastMessage: "Can we finalize the Reel delivery for Friday?",
    lastTimeIso: iso("2026-05-11", "10:20"),
    unreadCount: 2,
    online: true,
    replyHint: "Usually replies in 1 hour",
    messages: [
      {
        id: "m-1-1",
        role: "system",
        body: "Campaign invite sent",
        sentAt: iso("2026-05-09", "09:00"),
      },
      {
        id: "m-1-2",
        role: "business",
        body: "Hi Hana — excited to partner on Betu Cosmetics. Can you confirm your availability next week?",
        sentAt: iso("2026-05-09", "09:05"),
      },
      {
        id: "m-1-3",
        role: "influencer",
        body: "Hi! Yes, I’m open. I can deliver 1 TikTok + 2 stories within 5 days.",
        sentAt: iso("2026-05-09", "11:30"),
      },
      {
        id: "m-1-4",
        role: "business",
        body: "Perfect. Here’s the brief PDF for reference.",
        sentAt: iso("2026-05-10", "08:15"),
        attachmentLabel: "Brief_v1.pdf",
      },
      {
        id: "m-1-5",
        role: "influencer",
        body: "Can we finalize the Reel delivery for Friday?",
        sentAt: iso("2026-05-11", "10:20"),
      },
    ],
  },
  {
    id: "conv-2",
    influencerId: "inf-2",
    influencerName: "Biruk Alemu",
    imageUrl: "/images/influencers/influencer-2.jpg",
    campaignName: "Ride App Creator Push",
    niche: "Lifestyle",
    platform: "YouTube",
    lastMessage: "I can add a dedicated Shorts cut as well.",
    lastTimeIso: iso("2026-05-10", "16:40"),
    unreadCount: 0,
    online: false,
    replyHint: "Usually replies in 1 hour",
    messages: [
      {
        id: "m-2-1",
        role: "system",
        body: "Campaign invite sent",
        sentAt: iso("2026-05-08", "14:00"),
      },
      {
        id: "m-2-2",
        role: "business",
        body: "Hey Biruk — we’d love a creator-first walkthrough of the Ride app.",
        sentAt: iso("2026-05-08", "14:10"),
      },
      {
        id: "m-2-3",
        role: "influencer",
        body: "Sounds good. I can do a 6–8 min YouTube feature + cross-post Shorts.",
        sentAt: iso("2026-05-09", "10:00"),
      },
      {
        id: "m-2-4",
        role: "influencer",
        body: "I can add a dedicated Shorts cut as well.",
        sentAt: iso("2026-05-10", "16:40"),
      },
    ],
  },
  {
    id: "conv-3",
    influencerId: "inf-3",
    influencerName: "Meron Dawit",
    imageUrl: "/images/influencers/influencer-3.jpg",
    campaignName: "Coffee House Awareness",
    niche: "Fashion",
    platform: "Instagram",
    lastMessage: "Let’s lock the reel hook in the first 2 seconds.",
    lastTimeIso: iso("2026-05-11", "08:05"),
    unreadCount: 1,
    online: true,
    replyHint: "Online",
    messages: [
      {
        id: "m-3-1",
        role: "business",
        body: "Meron — for Coffee House we want warm, café lifestyle vibes.",
        sentAt: iso("2026-05-10", "12:00"),
      },
      {
        id: "m-3-2",
        role: "influencer",
        body: "Love it. I’ll scout two locations in Bole for B-roll.",
        sentAt: iso("2026-05-10", "15:22"),
      },
      {
        id: "m-3-3",
        role: "influencer",
        body: "Let’s lock the reel hook in the first 2 seconds.",
        sentAt: iso("2026-05-11", "08:05"),
      },
    ],
  },
  {
    id: "conv-4",
    influencerId: "inf-4",
    influencerName: "Saron Bekele",
    imageUrl: "/images/influencers/influencer-4.jpg",
    campaignName: "Abyssinia Restaurant Promo",
    niche: "Lifestyle",
    platform: "TikTok",
    lastMessage: "I’ll send a shot list tomorrow morning.",
    lastTimeIso: iso("2026-05-07", "19:50"),
    unreadCount: 0,
    online: false,
    replyHint: "Usually replies in 1 hour",
    messages: [
      {
        id: "m-4-1",
        role: "system",
        body: "Campaign invite sent",
        sentAt: iso("2026-05-06", "09:30"),
      },
      {
        id: "m-4-2",
        role: "influencer",
        body: "Abyssinia sounds fun — food + ambiance content is my sweet spot.",
        sentAt: iso("2026-05-06", "11:10"),
      },
      {
        id: "m-4-3",
        role: "business",
        body: "Great. We’re aiming for weekend traffic pushes.",
        sentAt: iso("2026-05-07", "09:00"),
      },
      {
        id: "m-4-4",
        role: "influencer",
        body: "I’ll send a shot list tomorrow morning.",
        sentAt: iso("2026-05-07", "19:50"),
      },
    ],
  },
];
