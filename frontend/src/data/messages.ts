import type { Message } from "@/types";

export const messages: Message[] = [
  {
    id: "msg-1",
    from: "Habesha Coffee Roasters",
    to: "Selam Beauty Studio",
    subject: "Content timeline confirmation",
    preview: "Can we finalize the Reel delivery for Friday?",
    unread: true,
    timestamp: "2026-05-07T10:20:00Z",
  },
  {
    id: "msg-2",
    from: "Addis Food Trails",
    to: "Abay Fashion House",
    subject: "Campaign brief feedback",
    preview: "I can include two teaser clips and one live stream.",
    unread: false,
    timestamp: "2026-05-06T15:42:00Z",
  },
  {
    id: "msg-3",
    from: "Adfluence Admin",
    to: "Sheba Tech Hub",
    subject: "Verification update",
    preview: "Your campaign assets are approved and ready to publish.",
    unread: false,
    timestamp: "2026-05-05T08:11:00Z",
  },
];
