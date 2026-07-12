import { CheckCircle, TrendingUp } from "lucide-react";

import { Container } from "@/components/common/container";

const features = [
  "Quality Influencers",
  "Data-Driven Results",
  "End-to-End Support",
  "Secure & Transparent",
];

export function WhyChooseSection() {
  return (
    <section className="bg-[color:var(--adfluence-surface)]/50 py-20">
      <Container>
        <div className="grid md:grid-cols-2 gap-12 items-center">
        <div>
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            Why Choose Adfluence?
          </h2>
          <p className="text-xl text-[#CBD5E1] mb-8">
            The most advanced influencer marketing platform built specifically
            for Ethiopia.
          </p>
          <div className="space-y-6">
            {features.map((feature, idx) => (
              <div key={feature} className="flex gap-4">
                <CheckCircle
                  className={`h-6 w-6 flex-shrink-0 ${
                    idx % 2 === 0
                      ? "text-[color:var(--adfluence-orange)]"
                      : "text-[color:var(--adfluence-purple-accent)]"
                  }`}
                />
                <div>
                  <h4 className="font-semibold mb-1">{feature}</h4>
                  <p className="text-sm text-[#CBD5E1]">
                    {feature === "Quality Influencers"
                      ? "Verified profiles with real engagement signals"
                      : feature === "Data-Driven Results"
                        ? "Track reach and ROI with transparent reporting"
                        : feature === "End-to-End Support"
                          ? "From briefs to approvals—everything in one place"
                          : "Clear contracts with secure collaboration flow"}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="adfluence-glass rounded-2xl p-6 border border-white/20">
          <div className="flex items-center justify-between pb-4 border-b border-white/10">
            <h4 className="font-semibold">Campaign Dashboard</h4>
            <span className="inline-flex items-center rounded-full bg-emerald-500/15 px-3 py-1 text-xs font-medium text-emerald-300 border border-emerald-500/30">
              Live
            </span>
          </div>
          <div className="grid grid-cols-2 gap-4 mt-4">
            <div className="bg-white/5 rounded-xl p-4">
              <div className="mb-1 text-2xl font-bold text-[color:var(--adfluence-orange)]">2.4M</div>
              <div className="text-xs text-[#CBD5E1]">Total Impressions</div>
            </div>
            <div className="bg-white/5 rounded-xl p-4">
              <div className="mb-1 text-2xl font-bold text-[color:var(--adfluence-purple-accent)]">8.2%</div>
              <div className="text-xs text-[#CBD5E1]">Engagement Rate</div>
            </div>
            <div className="bg-white/5 rounded-xl p-4">
              <div className="text-2xl font-bold text-emerald-300 mb-1">12K</div>
              <div className="text-xs text-[#CBD5E1]">Conversions</div>
            </div>
            <div className="bg-white/5 rounded-xl p-4">
              <div className="text-2xl font-bold text-amber-200 mb-1">3.8x</div>
              <div className="text-xs text-[#CBD5E1]">ROI</div>
            </div>
          </div>
          <div className="mt-4 h-32 bg-white/5 rounded-xl flex items-end p-4 gap-2">
            {[40, 65, 45, 80, 60, 90, 75].map((h, i) => (
              <div
                key={i}
                className="flex-1 rounded-t bg-[color:var(--adfluence-orange)]"
                style={{ height: `${h}%` }}
              />
            ))}
          </div>
          <div className="mt-4 flex items-center gap-2 text-xs text-[#CBD5E1]">
            <TrendingUp className="size-4 text-emerald-300" />
            Week-over-week performance trend
          </div>
        </div>
      </div>
      </Container>
    </section>
  );
}
