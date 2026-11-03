"use client";

import Link from "next/link";
import { useState } from "react";
import { ChevronDown } from "lucide-react";

import { useMarketingAuth } from "@/components/marketing/marketing-auth-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { formatEtb, PLATFORM_FEE_PERCENT, previewFees } from "@/lib/contracts";
import { parseAmount } from "@/lib/numbers";

const container = "mx-auto max-w-[1200px] px-5 sm:px-8";
const sectionTitle =
  "font-[family-name:var(--font-display)] text-[28px] leading-tight font-extrabold tracking-[-0.02em] text-white sm:text-[40px]";

const FREE = [
  "Creating an account, as a business or a creator",
  "Posting campaigns and getting matched with creators",
  "Inviting creators and chatting with them",
  "Verification and support from our team",
];

const FAQS = [
  {
    q: "When are the fees charged?",
    a: "Only when a contract is signed. The business pays the agreed amount plus 5%, which Adfluence holds safely. When the business approves the work, the creator receives the amount minus 5%.",
  },
  {
    q: "Do creators pay anything to join?",
    a: "No. Creators join for free. Their 5% comes out of the payment for a finished contract and is never charged up front.",
  },
  {
    q: "Are there monthly plans or subscriptions?",
    a: "No. There are no plans, setup fees or charges for browsing. You only pay when a real deal is agreed.",
  },
  {
    q: "What happens if the work isn't delivered?",
    a: "The money stays held by Adfluence until the business approves the work. If there is a problem, our team steps in to help both sides reach a fair outcome.",
  },
];

/** Live fee calculator: the same 5% + 5% arithmetic the contracts use. */
function FeeCalculator() {
  const [raw, setRaw] = useState("20,000");
  const amount = parseAmount(raw);
  const valid = amount !== undefined && Number.isFinite(amount) && amount > 0;
  const fees = previewFees(valid ? amount : 0);

  return (
    <div className="rounded-[24px] border border-white/10 bg-[color:var(--adfluence-card)] p-6 sm:p-8">
      <label htmlFor="fee-amount" className="text-sm font-semibold text-[#A78BFA]">
        Try it: the agreed price in birr
      </label>
      <div className="mt-3 flex items-center gap-3 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-4 focus-within:border-[color:var(--adfluence-orange)]">
        <span className="text-[15px] font-semibold text-[#94A3B8]">ETB</span>
        <Input
          id="fee-amount"
          inputMode="decimal"
          value={raw}
          onChange={(e) => setRaw(e.target.value)}
          className="h-14 border-0 bg-transparent px-0 text-2xl font-bold md:text-2xl text-white shadow-none focus-visible:ring-0 dark:bg-transparent"
        />
      </div>
      {!valid && raw.trim() ? <p className="mt-2 text-sm text-amber-200">Enter an amount, e.g. 15000 or 15k.</p> : null}
      <dl className="mt-6 space-y-4">
        <div className="flex items-baseline justify-between gap-4">
          <dt className="text-[15px] text-[#CBD5E1]">
            The business pays
            <span className="block text-xs text-[#94A3B8]">price + {PLATFORM_FEE_PERCENT}%</span>
          </dt>
          <dd className="text-xl font-bold text-white tabular-nums">{formatEtb(fees.totalBusinessPays)}</dd>
        </div>
        <div className="flex items-baseline justify-between gap-4">
          <dt className="text-[15px] text-[#CBD5E1]">
            The creator receives
            <span className="block text-xs text-[#94A3B8]">price − {PLATFORM_FEE_PERCENT}%</span>
          </dt>
          <dd className="text-xl font-bold text-white tabular-nums">{formatEtb(fees.influencerReceives)}</dd>
        </div>
        <div className="flex items-baseline justify-between gap-4 border-t border-white/10 pt-4">
          <dt className="text-[15px] text-[#CBD5E1]">Adfluence keeps</dt>
          <dd className="text-lg font-semibold text-[color:var(--adfluence-orange)] tabular-nums">
            {formatEtb(fees.businessFee + fees.influencerFee)}
          </dd>
        </div>
      </dl>
    </div>
  );
}

export function HowFeesWorkPage() {
  const { openSignup } = useMarketingAuth();

  return (
    <div className="bg-[color:var(--adfluence-bg)] text-white">
      <section className={`${container} grid items-center gap-10 py-12 sm:py-20 lg:grid-cols-[minmax(0,1fr)_minmax(0,480px)] lg:gap-16`}>
        <div>
          <h1 className="font-[family-name:var(--font-display)] text-[40px] leading-[1.04] font-extrabold tracking-[-0.03em] text-white sm:text-6xl">
            Simple fees, in birr.
          </h1>
          <p className="mt-5 max-w-[520px] text-[17px] leading-relaxed text-[#CBD5E1] sm:text-lg">
            Joining is free for everyone. When a contract is signed, the business pays {PLATFORM_FEE_PERCENT}% on top and
            the creator pays {PLATFORM_FEE_PERCENT}% out of the agreed price. Nothing else, ever.
          </p>
          <ul className="mt-7 space-y-3">
            {FREE.map((f) => (
              <li key={f} className="flex gap-3 text-[15px] text-[#CBD5E1]">
                <span aria-hidden className="font-bold text-[color:var(--adfluence-orange)]">
                  —
                </span>
                <span>
                  {f} <span className="text-[#94A3B8]">· free</span>
                </span>
              </li>
            ))}
          </ul>
        </div>
        <FeeCalculator />
      </section>

      <section className="bg-[color:var(--adfluence-surface)]">
        <div className={`${container} py-12 sm:py-20`}>
          <h2 className={sectionTitle}>Questions about fees</h2>
          <div className="mt-8 divide-y divide-white/10 rounded-[24px] border border-white/10 bg-[color:var(--adfluence-card)]">
            {FAQS.map((f) => (
              <details key={f.q} className="group px-6 py-5 sm:px-8">
                <summary className="flex cursor-pointer list-none items-center justify-between gap-4 text-base font-semibold text-white sm:text-lg [&::-webkit-details-marker]:hidden">
                  {f.q}
                  <ChevronDown aria-hidden className="size-5 shrink-0 text-[#94A3B8] transition-transform group-open:rotate-180" />
                </summary>
                <p className="mt-3 max-w-[760px] text-[15px] leading-relaxed text-[#CBD5E1]">{f.a}</p>
              </details>
            ))}
          </div>
        </div>
      </section>

      <section className={`${container} py-14 text-left sm:py-20 sm:text-center`}>
        <h2 className={sectionTitle}>Ready when you are</h2>
        <p className="mt-3 text-base text-[#CBD5E1] sm:text-lg">Free to join. You only pay when a deal is signed.</p>
        <div className="mt-7 flex flex-col gap-3 sm:flex-row sm:justify-center">
          <Button type="button" onClick={() => openSignup("business")} className="h-12 rounded-xl px-6 font-semibold">
            Post a campaign
          </Button>
          <Button asChild variant="secondary" className="h-12 rounded-xl border border-white/15 px-6 font-semibold">
            <Link href="/contact">Ask us a question</Link>
          </Button>
        </div>
      </section>
    </div>
  );
}
