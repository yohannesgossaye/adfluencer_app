"use client";

import { ArrowRight } from "lucide-react";

import { useMarketingAuth } from "@/components/marketing/marketing-auth-provider";
import { Button } from "@/components/ui/button";
import { formatEtb, previewFees } from "@/lib/contracts";
import { cn } from "@/lib/utils";

// The example deal uses the real fee arithmetic, so it can never drift from contracts.
const EXAMPLE = previewFees(20000);

function DealCard({
  step,
  title,
  tag,
  className,
  children,
}: {
  step: string;
  title: string;
  tag: React.ReactNode;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "w-full max-w-[440px] rounded-[20px] border border-white/10 bg-[color:var(--adfluence-card)] p-5 shadow-[0_18px_40px_rgba(2,8,18,0.55)] transition-transform duration-300 motion-safe:hover:rotate-0",
        className,
      )}
    >
      <div className="flex items-center justify-between gap-3">
        <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[#A78BFA]">
          {step} · {title}
        </p>
        {tag}
      </div>
      <div className="mt-3">{children}</div>
    </div>
  );
}

function Pill({ tone, children }: { tone: "success" | "info"; children: React.ReactNode }) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium",
        tone === "success"
          ? "border-emerald-400/40 bg-emerald-400/10 text-emerald-200"
          : "border-sky-400/40 bg-sky-400/10 text-sky-200",
      )}
    >
      <span aria-hidden className={cn("size-1.5 rounded-full", tone === "success" ? "bg-emerald-300" : "bg-sky-300")} />
      {children}
    </span>
  );
}

export function HeroDeal() {
  const { openSignup } = useMarketingAuth();

  return (
    <section className="mx-auto grid max-w-[1200px] items-center gap-12 px-5 pt-10 pb-14 sm:px-8 lg:grid-cols-[minmax(0,1fr)_minmax(0,460px)] lg:gap-16 lg:pt-16 lg:pb-24">
      <div>
        <h1 className="font-[family-name:var(--font-display)] text-[42px] leading-[1.02] font-extrabold tracking-[-0.03em] text-white sm:text-6xl lg:text-[72px] lg:leading-[0.98]">
          Find the creators Ethiopia actually listens to.
        </h1>
        <p lang="am" className="mt-5 font-[family-name:var(--font-amharic)] text-xl text-[#A78BFA] sm:text-2xl">
          ትክክለኛ ፈጣሪዎችን በቀላሉ ያግኙ።
        </p>
        <p className="mt-5 max-w-[560px] text-[17px] leading-relaxed text-[#CBD5E1] sm:text-[19px]">
          Post a brief. Get matched with vetted TikTok, Instagram and YouTube creators. Agree on one clear contract in
          birr, and pay only when the work is delivered.
        </p>
        <div className="mt-8 flex flex-col gap-3 sm:flex-row">
          <Button type="button" onClick={() => openSignup("business")} className="h-12 rounded-xl px-6 text-[15px] font-semibold">
            Post a campaign
          </Button>
          <Button
            type="button"
            variant="secondary"
            onClick={() => openSignup("influencer")}
            className="h-12 rounded-xl border border-white/15 px-6 text-[15px] font-semibold"
          >
            I&apos;m a creator <ArrowRight className="size-4" aria-hidden />
          </Button>
        </div>
        <p className="mt-4 text-sm text-[#94A3B8]">Free to join · 5% from each side, only on signed contracts</p>
      </div>

      {/* An illustration of one deal, not real users: labelled as such. */}
      <figure aria-label="Example of a deal on Adfluence" className="flex flex-col items-center">
        <div className="flex w-full flex-col items-center gap-3 lg:gap-0">
          <DealCard step="1" title="The brief" tag={<span className="text-xs text-[#94A3B8]">Example</span>} className="lg:rotate-[2.5deg]">
            <p className="font-[family-name:var(--font-display)] text-xl font-bold text-white">Coffee brand launch video</p>
            <p className="mt-1 text-sm text-[#CBD5E1]">TikTok · Food &amp; lifestyle · Addis Ababa · ETB 10,000–25,000</p>
          </DealCard>
          <DealCard step="2" title="The match" tag={<Pill tone="success">Strong fit</Pill>} className="lg:-mt-3 lg:-rotate-[1.5deg]">
            <div className="flex items-center gap-3">
              <span aria-hidden className="flex size-12 shrink-0 items-center justify-center rounded-full bg-[#4C1D95] font-semibold text-white">
                ST
              </span>
              <div>
                <p className="font-semibold text-white">Selam T.</p>
                <p className="text-[13px] text-[#CBD5E1]">Food creator · TikTok · verified by Adfluence</p>
              </div>
            </div>
            <p className="mt-3 text-sm text-[#CBD5E1]">Matches your niche, platform and budget</p>
          </DealCard>
          <DealCard step="3" title="The contract" tag={<Pill tone="info">Payment held safely</Pill>} className="lg:-mt-3 lg:rotate-[1.5deg]">
            <dl className="space-y-2 text-[15px]">
              <div className="flex justify-between">
                <dt className="text-[#CBD5E1]">Agreed amount</dt>
                <dd className="font-medium text-white">{formatEtb(EXAMPLE.amount)}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-[#CBD5E1]">Adfluence fee (5%)</dt>
                <dd className="font-medium text-white">{formatEtb(EXAMPLE.businessFee)}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-[#CBD5E1]">You pay</dt>
                <dd className="font-bold text-white">{formatEtb(EXAMPLE.totalBusinessPays)}</dd>
              </div>
            </dl>
            <p className="mt-3 text-[13px] text-[#94A3B8]">Released to the creator when you approve the video.</p>
          </DealCard>
        </div>
        <figcaption className="mt-5 text-xs text-[#94A3B8]">
          Illustration of how a deal looks. Names and amounts are examples.
        </figcaption>
      </figure>
    </section>
  );
}
