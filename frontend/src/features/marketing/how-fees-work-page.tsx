import Link from "next/link";
import { ArrowRight, BadgePercent, CircleDollarSign, Handshake, ShieldCheck } from "lucide-react";

import { Container } from "@/components/common/container";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { buildMarketingAuthHref } from "@/lib/marketing-auth-url";

const highlights = [
  {
    icon: CircleDollarSign,
    title: "Free to join",
    description:
      "Businesses and influencers can create accounts, browse the marketplace, and start conversations at no cost.",
  },
  {
    icon: BadgePercent,
    title: "5% per contract",
    description:
      "Adfluence charges a simple 5% platform fee only when a collaboration contract is completed.",
  },
  {
    icon: ShieldCheck,
    title: "No hidden fees",
    description:
      "No monthly subscriptions, no setup fees, and no charges for browsing creators or posting campaigns.",
  },
];

const included = [
  "Campaign creation and creator discovery",
  "Smart influencer recommendations",
  "In-app messaging and contract workflows",
  "Profile verification and marketplace visibility",
  "Support for Ethiopian brands and creators",
];

const faqs = [
  {
    question: "When is the 5% fee charged?",
    answer:
      "The platform fee applies when a contract is marked completed — not when you sign up, post a campaign, or send an invite.",
  },
  {
    question: "Do influencers pay to join?",
    answer:
      "No. Influencers join for free. The 5% fee is taken from the completed contract value, not from creator sign-up.",
  },
  {
    question: "Are there monthly plans?",
    answer:
      "No. Adfluence does not use subscription tiers. You only pay when real collaborations close on the platform.",
  },
  {
    question: "Who pays the platform fee?",
    answer:
      "The fee is calculated on the total contract value agreed between the business and influencer before payout.",
  },
  {
    question: "Is this built for Ethiopia?",
    answer:
      "Yes. Adfluence is designed for Ethiopian businesses and creators, with local pricing and campaign workflows in mind.",
  },
];

export function HowFeesWorkPage() {
  return (
    <div className="bg-[color:var(--adfluence-bg)] text-white">
      <section className="py-20">
        <Container>
          <div className="mx-auto max-w-3xl text-center">
            <p className="text-sm font-medium uppercase tracking-[0.2em] text-[color:var(--adfluence-orange)]">
              Transparent pricing
            </p>
            <h1 className="mt-4 text-3xl font-bold sm:text-4xl md:text-5xl">How Adfluence fees work</h1>
            <p className="mt-4 text-base text-[#CBD5E1] sm:text-lg">
              No subscriptions. No surprise charges. Just a simple 5% platform fee when a
              collaboration contract is completed.
            </p>
          </div>
        </Container>
      </section>

      <section className="pb-16">
        <Container>
          <div className="adfluence-glass mx-auto max-w-4xl rounded-3xl border border-[color:var(--adfluence-orange)]/30 p-8 text-center md:p-12">
            <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-[color:var(--adfluence-orange)]">
              <Handshake className="h-8 w-8 text-white" />
            </div>
            <p className="text-sm font-medium uppercase tracking-wide text-[#CBD5E1]">
              Platform fee
            </p>
            <p className="mt-2 text-4xl font-bold text-[color:var(--adfluence-orange)] sm:text-5xl md:text-6xl">
              5%
            </p>
            <p className="mt-3 text-xl text-white">per completed contract</p>
            <p className="mx-auto mt-4 max-w-2xl text-sm text-[#CBD5E1] md:text-base">
              You only pay when value is delivered. Browse creators, run campaigns, and negotiate
              deals for free — the fee applies once a contract is successfully completed.
            </p>
          </div>
        </Container>
      </section>

      <section className="pb-20">
        <Container>
          <div className="grid gap-6 md:grid-cols-3">
            {highlights.map((item) => (
              <article
                key={item.title}
                className="adfluence-glass rounded-2xl border border-white/15 p-6"
              >
                <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-[color:var(--adfluence-purple-mid)]">
                  <item.icon className="h-6 w-6 text-white" />
                </div>
                <h2 className="text-xl font-semibold text-white">{item.title}</h2>
                <p className="mt-2 text-sm leading-relaxed text-[#CBD5E1]">{item.description}</p>
              </article>
            ))}
          </div>
        </Container>
      </section>

      <section className="pb-20">
        <Container>
          <div className="grid gap-8 lg:grid-cols-2">
            <div className="adfluence-glass rounded-2xl border border-white/15 p-6 md:p-8">
              <h2 className="text-2xl font-bold text-white">Example calculation</h2>
              <p className="mt-2 text-sm text-[#CBD5E1]">
                If a business agrees to pay an influencer ETB 100,000 for a completed campaign
                contract:
              </p>
              <dl className="mt-6 space-y-4 text-sm">
                <div className="flex items-center justify-between border-b border-white/10 pb-3">
                  <dt className="text-[#CBD5E1]">Contract value</dt>
                  <dd className="font-semibold text-white">ETB 100,000</dd>
                </div>
                <div className="flex items-center justify-between border-b border-white/10 pb-3">
                  <dt className="text-[#CBD5E1]">Adfluence fee (5%)</dt>
                  <dd className="font-semibold text-[color:var(--adfluence-orange)]">ETB 5,000</dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-[#CBD5E1]">Influencer payout (before tax)</dt>
                  <dd className="font-semibold text-white">ETB 95,000</dd>
                </div>
              </dl>
              <p className="mt-4 text-xs text-[#94A3B8]">
                Demo illustration only. Final payout and tax handling will follow your contract
                terms and local regulations.
              </p>
            </div>

            <div className="adfluence-glass rounded-2xl border border-white/15 p-6 md:p-8">
              <h2 className="text-2xl font-bold text-white">Included at no extra cost</h2>
              <p className="mt-2 text-sm text-[#CBD5E1]">
                Everything you need to discover creators and run campaigns is free until a deal
                closes.
              </p>
              <ul className="mt-6 space-y-3 text-sm text-[#CBD5E1]">
                {included.map((item) => (
                  <li key={item} className="flex gap-3">
                    <span className="mt-1 text-[color:var(--adfluence-orange)]">•</span>
                    <span>{item}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </Container>
      </section>

      <section className="pb-16">
        <Container>
          <div className="rounded-2xl border border-white/15 bg-[color:var(--adfluence-surface)]/80 p-6 text-center text-[#CBD5E1] backdrop-blur-xl md:p-8">
            Influencers join for free and earn from brand collaborations. Businesses only pay the
            5% fee when a contract is completed.
          </div>
        </Container>
      </section>

      <section className="pb-20">
        <Container>
          <div className="mx-auto max-w-4xl">
            <h2 className="text-center text-3xl font-bold text-white">FAQ</h2>
            <div className="mt-8 space-y-3">
              {faqs.map((faq) => (
                <div
                  key={faq.question}
                  className="adfluence-glass rounded-2xl border border-white/15 p-5"
                >
                  <h3 className="text-lg font-semibold text-white">{faq.question}</h3>
                  <p className="mt-2 text-sm text-[#CBD5E1]">{faq.answer}</p>
                </div>
              ))}
            </div>
          </div>
        </Container>
      </section>

      <section className="relative overflow-hidden bg-[color:var(--adfluence-surface)] py-24">
        <div className="absolute top-1/2 left-1/2 h-96 w-96 -translate-x-1/2 -translate-y-1/2 rounded-full bg-[color:var(--adfluence-orange)]/8 blur-[120px]" />
        <Container className="relative text-center text-white">
          <h2 className="text-3xl font-bold sm:text-4xl md:text-5xl">Ready to launch your first campaign?</h2>
          <p className="mx-auto mt-4 max-w-2xl text-base text-[#CBD5E1] sm:text-lg">
            Join free, match with verified creators, and pay only when collaborations succeed.
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <GradientButton asChild className="h-12 rounded-xl px-6">
              <Link href={buildMarketingAuthHref("/pricing", "signup", "business")}>
                Get Started Free <ArrowRight className="size-5" />
              </Link>
            </GradientButton>
            <Button
              asChild
              variant="outline"
              className="h-12 rounded-xl border-white/25 bg-transparent px-6 text-white hover:bg-white/10"
            >
              <Link href={buildMarketingAuthHref("/pricing", "login")}>Talk to Us</Link>
            </Button>
          </div>
        </Container>
      </section>
    </div>
  );
}
