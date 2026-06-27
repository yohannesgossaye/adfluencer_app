import { CheckCircle2, LineChart } from "lucide-react";

import { Container } from "@/components/common/container";

const benefits = [
  "Trusted brands",
  "Secure payments",
  "Real campaign opportunities",
  "Transparent contracts",
  "Growth analytics",
];

export function InfluencerWhyLoveSection() {
  return (
    <section className="bg-[color:var(--adfluence-navy-dark)]/50 py-20">
      <Container>
        <div className="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <h2 className="text-3xl font-bold text-white sm:text-4xl">
              Why creators love Adfluence
            </h2>
            <div className="mt-6 space-y-4">
              {benefits.map((benefit) => (
                <div key={benefit} className="flex items-center gap-3">
                  <CheckCircle2 className="size-5 text-[#7C3AED]" />
                  <p className="text-[#CBD5E1]">{benefit}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="adfluence-glass rounded-2xl border border-white/20 p-6">
            <div className="flex items-center justify-between border-b border-white/10 pb-4">
              <h4 className="font-semibold text-white">Creator performance</h4>
              <span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-white">
                Updated
              </span>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-4">
              <div className="rounded-xl bg-white/5 p-4">
                <p className="text-xs text-[#CBD5E1]">Profile strength</p>
                <p className="mt-1 text-2xl font-semibold text-[#FF7A00]">89%</p>
              </div>
              <div className="rounded-xl bg-white/5 p-4">
                <p className="text-xs text-[#CBD5E1]">Invite acceptance</p>
                <p className="mt-1 text-2xl font-semibold text-[#7C3AED]">72%</p>
              </div>
            </div>

            <div className="mt-4 rounded-xl bg-white/5 p-4">
              <div className="flex items-center gap-2 text-xs text-[#CBD5E1]">
                <LineChart className="size-4 text-emerald-300" />
                Earnings graph
              </div>
              <div className="mt-3 flex h-24 items-end gap-2">
                {[25, 38, 47, 58, 70, 85, 92].map((h, i) => (
                  <div
                    key={i}
                    className="flex-1 rounded-t bg-gradient-to-t from-[#FF7A00] to-[#7C3AED]"
                    style={{ height: `${h}%` }}
                  />
                ))}
              </div>
            </div>

            <div className="mt-4 grid gap-2">
              {[
                "Timeline: 3 active collaborations",
                "Insight: +28% audience growth in 60 days",
                "Recommendation: Best posting window 7–9 PM",
              ].map((item) => (
                <div
                  key={item}
                  className="rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-xs text-[#CBD5E1]"
                >
                  {item}
                </div>
              ))}
            </div>
          </div>
        </div>
      </Container>
    </section>
  );
}

