import type { Metadata } from "next";

import { ContactPage } from "@/features/marketing/contact-page";

export const metadata: Metadata = {
  title: "Contact us · Adfluence",
  description: "Questions about campaigns, creators or fees? Send the Adfluence team a message; we reply by email.",
};

export default function MarketingContactPage() {
  return <ContactPage />;
}
