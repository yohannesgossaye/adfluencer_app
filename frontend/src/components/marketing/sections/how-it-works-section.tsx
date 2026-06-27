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
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED] flex items-center justify-center mb-6">
              <Target className="w-8 h-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold mb-6">For Businesses</h3>
            <div className="space-y-6">
              {["Create Campaign", "We Match You", "Get Results"].map((step, i) => (
                <div key={step} className="flex gap-4">
                  <div className="w-8 h-8 rounded-full bg-[#FF7A00]/20 flex items-center justify-center flex-shrink-0">
                    <span className="text-[#FF7A00] font-bold">{i + 1}</span>
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
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-[#7C3AED] to-[#FF7A00] flex items-center justify-center mb-6">
              <Sparkles className="w-8 h-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold mb-6">For Influencers</h3>
            <div className="space-y-6">
              {["Create Your Profile", "Get Matched", "Earn & Grow"].map((step, i) => (
                <div key={step} className="flex gap-4">
                  <div className="w-8 h-8 rounded-full bg-[#7C3AED]/20 flex items-center justify-center flex-shrink-0">
                    <span className="text-[#7C3AED] font-bold">{i + 1}</span>
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
