import { ArrowRight, CheckCircle2 } from "lucide-react";

import { Container } from "@/components/common/container";

const steps = [
  {
    title: "Create Campaign",
    description: "Set your goals, budget, platforms, and target audience.",
  },
  {
    title: "Get Matched",
    description: "Receive a shortlist of creators tailored to your brand.",
  },
  {
    title: "Collaborate",
    description: "Message, approve content, and align on deliverables.",
  },
  {
    title: "Launch & Track",
    description: "Monitor reach, engagement, conversions, and ROI in real time.",
  },
];

export function CampaignWorkflowSection() {
  return (
    <section className="py-20 bg-[color:var(--adfluence-navy-dark)]/50">
      <Container>
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-white sm:text-4xl">
            A clear campaign workflow
          </h2>
          <p className="mt-3 text-[#CBD5E1]">
            From brief to results — a modern flow that keeps teams aligned.
          </p>
        </div>

        <div className="grid gap-4 lg:grid-cols-4">
          {steps.map((s, i) => (
            <div key={s.title} className="relative">
              <div className="adfluence-glass rounded-2xl border border-white/15 p-6 h-full">
                <div className="flex items-center gap-3">
                  <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-white/10 text-sm font-semibold text-white">
                    {i + 1}
                  </span>
                  <h3 className="text-base font-semibold text-white">{s.title}</h3>
                </div>
                <p className="mt-3 text-sm leading-6 text-[#CBD5E1]">
                  {s.description}
                </p>
                <div className="mt-4 flex items-center gap-2 text-xs text-[#CBD5E1]">
                  <CheckCircle2 className="size-4 text-emerald-300" />
                  Structured and trackable
                </div>
              </div>
              {i < steps.length - 1 ? (
                <ArrowRight className="absolute -right-3 top-1/2 hidden size-5 -translate-y-1/2 text-white/30 lg:block" />
              ) : null}
            </div>
          ))}
        </div>
      </Container>
    </section>
  );
}

