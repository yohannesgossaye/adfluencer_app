import Image from "next/image";
import Link from "next/link";
import { ArrowRight, Heart, Play, TrendingUp, Zap } from "lucide-react";
import type { ReactNode } from "react";

import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function FloatingStat({
  title,
  value,
  icon,
  className,
}: {
  title: string;
  value: string;
  icon: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "absolute rounded-2xl border border-[rgba(255,255,255,0.22)] bg-[rgba(255,255,255,0.14)] px-4 py-3 text-white shadow-[0_16px_30px_rgba(6,17,31,0.28)] backdrop-blur-xl",
        className,
      )}
    >
      <div className="flex items-center gap-3">
        <span className="rounded-lg bg-white/20 p-2">{icon}</span>
        <div>
          <p className="text-xs text-slate-200">{title}</p>
          <p className="text-lg font-semibold">{value}</p>
        </div>
      </div>
    </div>
  );
}

export function HeroSection() {
  return (
    <section className="relative overflow-hidden bg-[color:var(--adfluence-navy)] text-white">
      <div className="absolute inset-y-0 right-0 hidden w-full md:block lg:w-[62%]">
        <Image
          src="/images/adfluence-hero.jpg"
          alt="Influencer creating content"
          fill
          priority
          className="object-cover object-[78%_center]"
        />
      </div>
      <div className="absolute inset-0 bg-[linear-gradient(90deg,rgba(11,26,43,1)_0%,rgba(11,26,43,0.95)_38%,rgba(11,26,43,0.75)_58%,rgba(11,26,43,0.28)_78%,rgba(11,26,43,0)_100%)]" />

      <div className="relative z-20 flex min-h-[calc(100vh-72px)] items-center px-4 sm:px-6 md:px-10 lg:px-16 xl:px-20">
        <div className="max-w-[560px] text-left">
          <p className="inline-flex rounded-full border border-white/14 bg-white/10 px-4 py-2 text-sm tracking-wide text-slate-200 backdrop-blur-sm">
            ETHIOPIA'S CREATOR MARKETING PLATFORM
          </p>
          <h1 className="mt-8 text-5xl font-bold leading-[1.03] sm:text-6xl lg:text-7xl">
            Connect with Influencers.
            <br />
            Grow Your Brand.
            <br />
            <span className="bg-gradient-to-r from-[#FF7A00] to-[#7C3AED] bg-clip-text text-transparent">
              Anywhere.
            </span>
          </h1>
          <p className="mt-6 text-lg leading-8 text-[#CBD5E1] lg:text-xl lg:leading-9">
            Adfluence helps businesses of all sizes discover the right creators,
            run successful campaigns, and achieve real results.
          </p>
          <div className="mt-8 flex flex-wrap items-center justify-start gap-4">
            <GradientButton
              asChild
              className="h-12 justify-center rounded-xl px-6 text-base font-medium sm:text-lg"
            >
              <Link href="/auth/signup?role=business">
                Start a Campaign <ArrowRight className="size-5" />
              </Link>
            </GradientButton>
            <Button
              asChild
              variant="outline"
              className="h-12 justify-center rounded-xl border-white/15 bg-transparent px-6 text-base font-medium text-white hover:bg-white/10 sm:text-lg"
            >
              <Link href="/auth/signup?role=influencer">
                <Play className="size-5" />
                Join as Influencer
              </Link>
            </Button>
          </div>
          <p className="mt-8 text-sm text-slate-300">
            Built for brands, creators, and agencies in Ethiopia.
          </p>
        </div>

        <FloatingStat
          title="Likes"
          value="12.8K"
          icon={<Heart className="size-4 text-rose-300" />}
          className="right-6 top-[55%] hidden lg:block"
        />
        <FloatingStat
          title="Reach"
          value="45.3K"
          icon={<TrendingUp className="size-4 text-violet-200" />}
          className="bottom-16 right-20 hidden lg:block"
        />
        <FloatingStat
          title="Match Score"
          value="92%"
          icon={<Zap className="size-4 text-amber-200" />}
          className="top-[72%] right-[30%] hidden lg:block"
        />
      </div>
    </section>
  );
}
