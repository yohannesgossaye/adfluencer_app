import { ArrowRight, BarChart3, Sparkles, TrendingUp, Users } from "lucide-react";
import Link from "next/link";

import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { buildMarketingAuthHref } from "@/lib/marketing-auth-url";

function MiniStat({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
      <p className="text-xl font-semibold text-white">{value}</p>
      <p className="text-xs text-[#CBD5E1]">{label}</p>
    </div>
  );
}

function CreatorMiniCard({ name, niche }: { name: string; niche: string }) {
  return (
    <div className="flex items-center justify-between rounded-2xl border border-white/10 bg-white/5 p-4">
      <div>
        <p className="text-sm font-semibold text-white">{name}</p>
        <p className="text-xs text-[#CBD5E1]">{niche}</p>
      </div>
      <span className="rounded-full bg-white/10 px-3 py-1 text-xs text-white">
        Verified
      </span>
    </div>
  );
}

export function BusinessHeroSection() {
  return (
    <section className="relative flex min-h-[calc(100vh-72px)] items-center overflow-hidden bg-[color:var(--adfluence-bg)] py-10 md:py-12 lg:py-16">
      <div className="absolute -left-24 top-8 h-72 w-72 rounded-full bg-[#0B2240]/40 blur-[110px]" />
      <div className="absolute right-[-12%] top-[15%] h-80 w-80 rounded-full bg-[color:var(--adfluence-orange)]/8 blur-[120px]" />
      <div className="relative w-full px-4 sm:px-6 md:px-10 lg:px-16 xl:px-20">
        <div className="grid items-center gap-10 lg:grid-cols-2">
          <div className="max-w-[560px]">
            <p className="inline-flex rounded-full border border-white/14 bg-white/10 px-4 py-2 text-xs tracking-wide text-[#CBD5E1] backdrop-blur-sm">
              BUILT FOR ETHIOPIAN BRANDS
            </p>
            <h1 className="mt-6 text-4xl font-bold leading-[1.05] text-white sm:text-5xl lg:text-6xl">
              Launch influencer campaigns that actually{" "}
              <span className="text-[color:var(--adfluence-orange)]">convert.</span>
            </h1>
            <p className="mt-5 text-base leading-7 text-[#CBD5E1] sm:text-lg">
              Discover creators, manage campaigns, track ROI, and collaborate with
              influencers — all in one focused workspace.
            </p>

            <div className="mt-8 flex flex-wrap gap-3">
              <GradientButton asChild className="h-12 rounded-xl px-6">
                <Link href={buildMarketingAuthHref("/for-businesses", "signup", "business")}>
                  Start a Campaign <ArrowRight className="size-5" />
                </Link>
              </GradientButton>
              <Button
                asChild
                variant="outline"
                className="h-12 rounded-xl border-white/15 bg-transparent px-6 text-white hover:bg-white/10"
              >
                <Link href={buildMarketingAuthHref("/for-businesses", "login")}>
                  Book a Demo
                </Link>
              </Button>
            </div>
          </div>

          <div className="adfluence-glass rounded-2xl border border-white/20 p-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-[color:var(--adfluence-purple)] text-white">
                  <BarChart3 className="size-5" />
                </span>
                <div>
                  <p className="text-sm font-semibold text-white">
                    Business Dashboard
                  </p>
                  <p className="text-xs text-[#CBD5E1]">Campaign overview</p>
                </div>
              </div>
              <span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-white">
                Live
              </span>
            </div>

            <div className="mt-6 grid gap-3 sm:grid-cols-3">
              <MiniStat label="Active Campaigns" value="8" />
              <MiniStat label="Matched Creators" value="24" />
              <MiniStat label="Avg. ROI" value="3.8x" />
            </div>

            <div className="mt-6 grid gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <div className="flex items-center gap-2 text-xs text-[#CBD5E1]">
                  <TrendingUp className="size-4 text-emerald-300" />
                  Weekly performance
                </div>
                <div className="mt-3 flex h-24 items-end gap-2">
                  {[40, 65, 45, 80, 60, 90, 75].map((h, i) => (
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
                  <Users className="size-4 text-[color:var(--adfluence-purple-accent)]" />
                  Creator shortlist
                </div>
                <CreatorMiniCard name="Addis Food Trails" niche="Food • TikTok" />
                <CreatorMiniCard name="Selam Beauty Studio" niche="Beauty • Instagram" />
              </div>
            </div>

            <div className="mt-6 flex flex-wrap items-center gap-2 text-xs text-[#CBD5E1]">
              <Sparkles className="size-4 text-[color:var(--adfluence-orange)]" />
              AI recommendations • Secure contracts • Transparent reporting
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
