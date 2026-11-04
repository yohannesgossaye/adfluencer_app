import type { Metadata } from "next";

import { ForInfluencersMarketingPage } from "@/features/influencer";

export const metadata: Metadata = {
  title: "For creators · Adfluence",
  description: "Join free, get approved, and receive real campaign briefs from Ethiopian businesses, with the price agreed and the payment secured first.",
};

export default function ForInfluencersRoutePage() {
  return <ForInfluencersMarketingPage />;
}
