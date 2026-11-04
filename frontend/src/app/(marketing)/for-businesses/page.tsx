import type { Metadata } from "next";

import { ForBusinessesMarketingPage } from "@/features/business";

export const metadata: Metadata = {
  title: "For businesses · Adfluence",
  description: "Find checked Ethiopian creators, see why each one fits your campaign, and agree everything in one contract in birr.",
};

export default function ForBusinessesRoutePage() {
  return <ForBusinessesMarketingPage />;
}
