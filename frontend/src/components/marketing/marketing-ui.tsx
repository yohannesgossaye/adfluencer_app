"use client";

import Link from "next/link";
import type { ReactNode } from "react";
import { ChevronDown } from "lucide-react";

import { useMarketingAuth } from "@/components/marketing/marketing-auth-provider";
import { Button } from "@/components/ui/button";
import type { MarketingAuthRole } from "@/lib/marketing-auth-url";
import { cn } from "@/lib/utils";

import { mContainer, mTitle } from "@/components/marketing/marketing-styles";

/** Shared building blocks for the public pages, so they share one rhythm. */

export function PageHero({ title, lead, children }: { title: string; lead: string; children?: ReactNode }) {
  return (
    <section className={cn(mContainer, "py-12 sm:py-20")}>
      <h1 className="max-w-[860px] font-[family-name:var(--font-display)] text-[40px] leading-[1.04] font-extrabold tracking-[-0.03em] text-white sm:text-6xl">
        {title}
      </h1>
      <p className="mt-5 max-w-[640px] text-[17px] leading-relaxed text-[#CBD5E1] sm:text-lg">{lead}</p>
      {children}
    </section>
  );
}

export function SignupButtons({ role, primary }: { role: MarketingAuthRole; primary: string }) {
  const { openSignup } = useMarketingAuth();
  return (
    <div className="mt-8 flex flex-col gap-3 sm:flex-row">
      <Button type="button" onClick={() => openSignup(role)} className="h-12 rounded-xl px-6 text-[15px] font-semibold">
        {primary}
      </Button>
      <Button asChild variant="secondary" className="h-12 rounded-xl border border-white/15 px-6 text-[15px] font-semibold">
        <Link href="/pricing">See the fees</Link>
      </Button>
    </div>
  );
}

export function Faq({ title, items }: { title: string; items: Array<{ q: string; a: string }> }) {
  return (
    <section className="bg-[color:var(--adfluence-surface)]">
      <div className={cn(mContainer, "py-12 sm:py-20")}>
        <h2 className={mTitle}>{title}</h2>
        <div className="mt-8 divide-y divide-white/10 rounded-[24px] border border-white/10 bg-[color:var(--adfluence-card)]">
          {items.map((f) => (
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
  );
}

export function ClosingCta({ title, lead, role, primary }: { title: string; lead: string; role: MarketingAuthRole; primary: string }) {
  const { openSignup } = useMarketingAuth();
  return (
    <section className={cn(mContainer, "py-14 text-left sm:py-20 sm:text-center")}>
      <h2 className={mTitle}>{title}</h2>
      <p className="mt-3 text-base text-[#CBD5E1] sm:text-lg">{lead}</p>
      <div className="mt-7 flex flex-col gap-3 sm:flex-row sm:justify-center">
        <Button type="button" onClick={() => openSignup(role)} className="h-12 rounded-xl px-6 font-semibold">
          {primary}
        </Button>
        <Button asChild variant="secondary" className="h-12 rounded-xl border border-white/15 px-6 font-semibold">
          <Link href="/contact">Ask us a question</Link>
        </Button>
      </div>
    </section>
  );
}

/** A numbered list of short points, one per row (not a grid of identical cards). */
export function NumberedPoints({ items }: { items: Array<{ title: string; text: string }> }) {
  return (
    <ol className="mt-8 divide-y divide-white/10 border-y border-white/10">
      {items.map((p, i) => (
        <li key={p.title} className="grid gap-2 py-6 sm:grid-cols-[80px_minmax(0,280px)_minmax(0,1fr)] sm:gap-6 sm:py-7">
          <span className="text-[15px] font-bold text-[color:var(--adfluence-orange)]">{String(i + 1).padStart(2, "0")}</span>
          <h3 className="font-[family-name:var(--font-display)] text-xl font-bold text-white">{p.title}</h3>
          <p className="text-[15px] leading-relaxed text-[#CBD5E1]">{p.text}</p>
        </li>
      ))}
    </ol>
  );
}
