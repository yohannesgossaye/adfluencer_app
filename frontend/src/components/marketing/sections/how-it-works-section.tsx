import { BadgeCheck, BriefcaseBusiness, Sparkles, Target } from "lucide-react";

import { Container } from "@/components/common/container";

const flows = [
  {
    title: "For Businesses",
    icon: BriefcaseBusiness,
    steps: ["Create Campaign", "We Match You", "Get Results"],
  },
  {
    title: "For Influencers",
    icon: Sparkles,
    steps: ["Create Your Profile", "Get Matched", "Earn & Grow"],
  },
];

export function HowItWorksSection() {
  return (
    <section id="how-it-works" className="py-20">
      <Container>
        <div className="text-center mb-16">
          <h2 className="text-4xl md:text-5xl font-bold mb-4">
            How Adfluence Works
          </h2>
          <p className="text-xl text-[#CBD5E1]">Simple, powerful, effective</p>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          <div className="adfluence-glass rounded-2xl p-8 border border-white/20">
            <div className="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-[color:var(--adfluence-orange)]">
              <Target className="h-8 w-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold mb-6">For Businesses</h3>
            <div className="space-y-6">
              {["Create Campaign", "We Match You", "Get Results"].map((step, i) => (
                <div key={step} className="flex gap-4">
                  <div className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-[color:var(--adfluence-orange)]/20">
                    <span className="font-bold text-[color:var(--adfluence-orange)]">{i + 1}</span>
                  </div>
                  <div>
                    <h4 className="font-semibold mb-1">{step}</h4>
                    <p className="text-sm text-[#CBD5E1]">
                      {i === 0
                        ? "Define your goals, budget, and audience"
                        : i === 1
                          ? "Get creator recommendations instantly"
                          : "Launch and track real performance"}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="adfluence-glass rounded-2xl p-8 border border-white/20">
            <div className="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-[color:var(--adfluence-purple-mid)]">
              <Sparkles className="h-8 w-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold mb-6">For Influencers</h3>
            <div className="space-y-6">
              {["Create Your Profile", "Get Matched", "Earn & Grow"].map((step, i) => (
                <div key={step} className="flex gap-4">
                  <div className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-[color:var(--adfluence-purple-mid)]/20">
                    <span className="font-bold text-[color:var(--adfluence-purple-accent)]">{i + 1}</span>
                  </div>
                  <div>
                    <h4 className="font-semibold mb-1">{step}</h4>
                    <p className="text-sm text-[#CBD5E1]">
                      {i === 0
                        ? "Showcase your niche and audience"
                        : i === 1
                          ? "Receive invites from top brands"
                          : "Collaborate, deliver, and get paid"}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </Container>
    </section>
  );
}
