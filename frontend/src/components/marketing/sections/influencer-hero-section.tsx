import { ArrowRight, BarChart3, DollarSign, Sparkles, Users2 } from "lucide-react";
import Link from "next/link";

import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { buildMarketingAuthHref } from "@/lib/marketing-auth-url";

function Metric({
  label,
  value,
  tone,
}: {
  label: string;
  value: string;
  tone: "orange" | "purple" | "green";
}) {
  const toneClass =
    tone === "orange"
      ? "text-[color:var(--adfluence-orange)]"
      : tone === "purple"
        ? "text-[color:var(--adfluence-purple-accent)]"
        : "text-emerald-300";

  return (
    <div className="rounded-xl border border-white/10 bg-white/5 p-4">
      <p className={`text-2xl font-bold ${toneClass}`}>{value}</p>
      <p className="text-xs text-[#CBD5E1]">{label}</p>
    </div>
  );
}

export function InfluencerHeroSection() {
  return (
    <section className="relative flex min-h-[calc(100vh-72px)] items-center overflow-hidden bg-[color:var(--adfluence-bg)] py-10 md:py-12 lg:py-16">
      <div className="absolute -left-20 top-12 h-72 w-72 rounded-full bg-[#0B2240]/40 blur-[110px]" />
      <div className="absolute right-[-14%] top-[18%] h-80 w-80 rounded-full bg-[color:var(--adfluence-purple-mid)]/12 blur-[120px]" />

      <div className="relative w-full px-4 sm:px-6 md:px-10 lg:px-16 xl:px-20">
        <div className="grid items-center gap-10 lg:grid-cols-2">
          <div className="max-w-[560px] text-left">
            <p className="inline-flex rounded-full border border-white/14 bg-white/10 px-4 py-2 text-xs tracking-wide text-[#CBD5E1] backdrop-blur-sm">
              BUILT FOR ETHIOPIAN CREATORS
            </p>
            <h1 className="mt-6 text-4xl font-bold leading-[1.05] text-white sm:text-5xl lg:text-6xl">
              Turn your audience into{" "}
              <span className="text-[color:var(--adfluence-orange)]">real income.</span>
            </h1>
            <p className="mt-5 text-base leading-7 text-[#CBD5E1] sm:text-lg">
              Build your profile, receive campaign invites, collaborate with
              trusted brands, and grow your creator career with transparent
              payouts.
            </p>

            <div className="mt-8 flex flex-wrap gap-3">
              <GradientButton asChild className="h-12 rounded-xl px-6">
                <Link href={buildMarketingAuthHref("/for-influencers", "signup", "influencer")}>
                  Join as Influencer <ArrowRight className="size-5" />
                </Link>
              </GradientButton>
              <Button
                asChild
                variant="outline"
                className="h-12 rounded-xl border-white/15 bg-transparent px-6 text-white hover:bg-white/10"
              >
                <Link href="/business/campaigns">Explore Campaigns</Link>
              </Button>
            </div>
          </div>

          <div className="adfluence-glass rounded-2xl border border-white/20 p-6">
            <div className="flex items-center justify-between border-b border-white/10 pb-4">
              <div className="flex items-center gap-3">
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-[color:var(--adfluence-purple)] text-white">
                  <Sparkles className="size-5" />
                </span>
                <div>
                  <p className="text-sm font-semibold text-white">
                    Creator Dashboard
                  </p>
                  <p className="text-xs text-[#CBD5E1]">Growth and earnings</p>
                </div>
              </div>
              <span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-white">
                Active
              </span>
            </div>

            <div className="mt-6 grid gap-3 sm:grid-cols-3">
              <Metric label="Campaign Invites" value="12" tone="purple" />
              <Metric label="Monthly Earnings" value="ETB 48K" tone="orange" />
              <Metric label="Completion Rate" value="94%" tone="green" />
            </div>

            <div className="mt-6 grid gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <div className="flex items-center gap-2 text-xs text-[#CBD5E1]">
                  <BarChart3 className="size-4 text-[color:var(--adfluence-purple-accent)]" />
                  Engagement analytics
                </div>
                <div className="mt-3 flex h-24 items-end gap-2">
                  {[35, 50, 62, 45, 70, 85, 78].map((h, i) => (
                    <div
                      key={i}
                      className="flex-1 rounded-t bg-[color:var(--adfluence-orange)]"
                      style={{ height: `${h}%` }}
                    />
                  ))}
                </div>
              </div>
              <div className="space-y-3">
                <div className="flex items-center gap-2 text-xs text-[#CBD5E1]">
                  <Users2 className="size-4 text-emerald-300" />
                  Collaboration requests
                </div>
                {[
                  "Habesha Coffee • Product Reel",
                  "Abay Fashion • Launch Teaser",
                  "Ride Ethiopia • Creator Promo",
                ].map((invite) => (
                  <div
                    key={invite}
                    className="rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-sm text-white"
                  >
                    {invite}
                  </div>
                ))}
              </div>
            </div>

            <div className="mt-6 flex flex-wrap items-center gap-2 text-xs text-[#CBD5E1]">
              <DollarSign className="size-4 text-[color:var(--adfluence-orange)]" />
              Secure payouts • Transparent contracts • Verified brand system
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
