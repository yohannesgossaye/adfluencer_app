"use client";

import { Mail, MessageCircle, Phone } from "lucide-react";

import { useMarketingAuth } from "@/components/marketing/marketing-auth-provider";
import { Button } from "@/components/ui/button";
import { formatEtb, previewFees } from "@/lib/contracts";
import { contactLinks } from "@/lib/site-contact";
import { cn } from "@/lib/utils";

const container = "mx-auto max-w-[1200px] px-5 sm:px-8";
const sectionTitle =
  "font-[family-name:var(--font-display)] text-[30px] leading-tight font-extrabold tracking-[-0.02em] text-white sm:text-[44px]";

const STEPS = [
  {
    title: "Post a brief",
    text: "Say what you're launching, where, on which platform and your budget in birr. It stays a draft until you publish.",
  },
  {
    title: "Get matched",
    text: "We rank approved creators by niche, platform, budget and audience. You see why each one fits.",
  },
  {
    title: "Agree in one contract",
    text: "Deliverables, dates and the price, with the 5% fee shown upfront. The creator accepts or declines.",
  },
  {
    title: "Paid on approval",
    text: "Your payment is held by Adfluence and released to the creator when you approve the work.",
  },
];

export function HowADealWorks() {
  return (
    <section id="how-it-works" className={cn(container, "scroll-mt-24 py-14 sm:py-24")}>
      <h2 className={sectionTitle}>How a deal works</h2>
      <ol className="mt-8 grid gap-8 sm:mt-12 sm:grid-cols-2 sm:gap-6 lg:grid-cols-4">
        {STEPS.map((s, i) => (
          <li key={s.title} className="flex gap-4 sm:block sm:border-t-[3px] sm:border-[color:var(--adfluence-action)] sm:pt-6">
            <span className="text-[15px] font-bold text-[color:var(--adfluence-orange)]">{String(i + 1).padStart(2, "0")}</span>
            <div>
              <h3 className="font-[family-name:var(--font-display)] text-lg font-bold text-white sm:mt-3 sm:text-[22px]">
                {s.title}
              </h3>
              <p className="mt-1 text-[15px] leading-relaxed text-[#CBD5E1] sm:mt-3">{s.text}</p>
            </div>
          </li>
        ))}
      </ol>
    </section>
  );
}

function Side({
  title,
  points,
  action,
  className,
}: {
  title: string;
  points: string[];
  action: React.ReactNode;
  className?: string;
}) {
  return (
    <div className={cn("flex flex-col rounded-[24px] p-6 sm:rounded-[28px] sm:p-12", className)}>
      <h2 className="font-[family-name:var(--font-display)] text-2xl font-extrabold tracking-[-0.01em] text-white sm:text-[34px]">
        {title}
      </h2>
      <ul className="mt-5 flex-1 space-y-3">
        {points.map((p) => (
          <li key={p} className="flex gap-3 text-[15px] leading-relaxed text-[#CBD5E1] sm:text-base">
            <span aria-hidden className="font-bold text-[color:var(--adfluence-orange)]">
              —
            </span>
            {p}
          </li>
        ))}
      </ul>
      <div className="mt-7">{action}</div>
    </div>
  );
}

export function TwoSides() {
  const { openSignup } = useMarketingAuth();
  return (
    <section className={cn(container, "grid gap-4 pb-14 sm:pb-24 md:grid-cols-2 md:gap-6")}>
      <Side
        title="For businesses"
        className="bg-[#4C1D95]/35"
        points={[
          "Only verified creators, with their real followers and prices",
          "See why each creator fits your campaign",
          "One contract and one payment, in birr",
          "Your money is released only when you approve the work",
        ]}
        action={
          <Button type="button" onClick={() => openSignup("business")} className="h-12 w-full rounded-xl px-6 font-semibold sm:w-auto">
            Post a campaign
          </Button>
        }
      />
      <Side
        title="For creators"
        className="bg-[color:var(--adfluence-card)]"
        points={[
          "Brands come to you with real campaign briefs",
          "Know the price and the deliverables before you say yes",
          "Payment is secured before you start",
          "Free to join, 5% only on deals you sign",
        ]}
        action={
          <Button
            type="button"
            variant="secondary"
            onClick={() => openSignup("influencer")}
            className="h-12 w-full rounded-xl border border-white/15 px-6 font-semibold sm:w-auto"
          >
            Create my creator profile
          </Button>
        }
      />
    </section>
  );
}

const FEES = previewFees(100000);

export function HonestFees() {
  return (
    <section className="bg-[color:var(--adfluence-surface)]">
      <div className={cn(container, "grid items-center gap-8 py-12 sm:py-20 lg:grid-cols-2 lg:gap-20")}>
        <div>
          <h2 className={sectionTitle}>Honest fees, in birr</h2>
          <p className="mt-4 max-w-[520px] text-base leading-relaxed text-[#CBD5E1] sm:text-lg">
            Joining, posting campaigns and chatting are free. When a contract is signed, the business pays 5% on top and
            the creator pays 5% out of the agreed amount. Nothing else.
          </p>
        </div>
        <div className="rounded-[20px] bg-[color:var(--adfluence-card)] p-6 sm:p-8">
          <p className="text-sm font-semibold text-[#A78BFA]">Example: a {formatEtb(FEES.amount)} contract</p>
          <dl className="mt-4 space-y-3 text-base sm:text-[17px]">
            {[
              ["Business pays", FEES.totalBusinessPays],
              ["Creator receives", FEES.influencerReceives],
              ["Adfluence keeps", FEES.businessFee + FEES.influencerFee],
            ].map(([label, value]) => (
              <div key={label} className="flex justify-between gap-4">
                <dt className="text-[#CBD5E1]">{label}</dt>
                <dd className="font-bold text-white">{formatEtb(value as number)}</dd>
              </div>
            ))}
          </dl>
        </div>
      </div>
    </section>
  );
}

export function TalkToUs() {
  const contact = "h-12 w-full rounded-xl border border-white/15 px-5 font-semibold sm:w-auto";
  return (
    <section id="talk-to-us" className={cn(container, "scroll-mt-24 py-14 text-left sm:py-24 sm:text-center")}>
      <h2 className={sectionTitle}>Talk to a person, not a form</h2>
      <p className="mt-3 text-base text-[#CBD5E1] sm:text-lg">Our team in Addis Ababa answers in Amharic or English.</p>
      <div className="mt-7 flex flex-col gap-3 sm:flex-row sm:justify-center">
        <Button asChild variant="secondary" className={contact}>
          <a href={contactLinks.email}>
            <Mail className="size-4" aria-hidden /> Email us
          </a>
        </Button>
        <Button asChild variant="secondary" className={contact}>
          <a href={contactLinks.phone}>
            <Phone className="size-4" aria-hidden /> Call us
          </a>
        </Button>
        <Button asChild variant="secondary" className={contact}>
          <a href={contactLinks.telegram} target="_blank" rel="noopener noreferrer">
            <MessageCircle className="size-4" aria-hidden /> Message on Telegram
          </a>
        </Button>
      </div>
    </section>
  );
}
