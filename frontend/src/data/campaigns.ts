import type { Campaign } from "@/types";

export const campaigns: Campaign[] = [
  {
    id: "cmp-1",
    businessName: "Habesha Coffee Roasters",
    title: "Summer Cold Brew Launch",
    budgetEtb: 250000,
    status: "active",
    platform: "Instagram",
  },
  {
    id: "cmp-2",
    businessName: "Abay Fashion House",
    title: "Enkutatash Collection Teaser",
    budgetEtb: 180000,
    status: "pending",
    platform: "TikTok",
  },
  {
    id: "cmp-3",
    businessName: "Sheba Tech Hub",
    title: "Campus Creator Challenge",
    budgetEtb: 300000,
    status: "draft",
    platform: "YouTube",
  },
  {
    id: "cmp-4",
    businessName: "Dallol Travel Co.",
    title: "Weekend Escape Promo",
    budgetEtb: 120000,
    status: "completed",
    platform: "Instagram",
  },
];
