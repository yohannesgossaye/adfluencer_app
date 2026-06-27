import {
  BadgeCheck,
  CreditCard,
  FileText,
  MessageSquare,
  Search,
  TrendingUp,
} from "lucide-react";

import { Container } from "@/components/common/container";

const highlights = [
  { title: "Smart Matching", icon: Search },
  { title: "Secure Contracts", icon: FileText },
  { title: "Fast Payouts", icon: CreditCard },
  { title: "Creator Analytics", icon: TrendingUp },
  { title: "Campaign Messaging", icon: MessageSquare },
  { title: "Verified Brand System", icon: BadgeCheck },
];

export function InfluencerFeatureHighlightsSection() {
  return (
    <section className="py-20">
      <Container>
        <div className="mb-12 text-center">
          <h2 className="text-3xl font-bold text-white sm:text-4xl">
            Feature highlights for creators
          </h2>
          <p className="mt-3 text-[#CBD5E1]">
            Creator-first tools designed for trust, speed, and growth.
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {highlights.map((feature) => (
            <article
              key={feature.title}
              className="adfluence-glass rounded-2xl border border-white/15 p-6"
            >
              <span className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED] text-white">
                <feature.icon className="size-5" />
              </span>
              <h3 className="mt-4 text-lg font-semibold text-white">
                {feature.title}
              </h3>
              <p className="mt-2 text-sm text-[#CBD5E1]">
                Built to make your campaign workflow smoother and more reliable.
              </p>
            </article>
          ))}
        </div>
      </Container>
    </section>
  );
}

