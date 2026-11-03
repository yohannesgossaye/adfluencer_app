import type { Metadata } from "next";

import { HowFeesWorkPage } from "@/features/marketing/how-fees-work-page";

export const metadata: Metadata = {
  title: "Fees · Adfluence",
  description: "Free to join. On a signed contract the business pays 5% on top and the creator 5% out of the price. Try the fee calculator.",
};

export default function PricingPage() {
  return <HowFeesWorkPage />;
}
