import { CheckCircle2, Sparkles } from "lucide-react";

import { Container } from "@/components/common/container";

const steps = [
  {
    title: "Create Profile",
    description:
      "Showcase your niche, audience, and best content in one verified profile.",
  },
  {
    title: "Receive Campaign Invites",
    description:
      "Get matched with relevant opportunities from trusted brands across Ethiopia.",
  },
  {
    title: "Collaborate & Earn",
    description:
      "Deliver content, track progress, and receive payouts with transparent terms.",
  },
];

export function InfluencerHowItWorksSection() {
  return (
    <section className="py-20">
      <Container>
        <div className="mb-12 text-center">
          <h2 className="text-3xl font-bold text-white sm:text-4xl">
            How it works for creators
          </h2>
          <p className="mt-3 text-[#CBD5E1]">
            A simple creator workflow from setup to earnings.
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-3">
          {steps.map((step, index) => (
            <article
              key={step.title}
              className="adfluence-glass rounded-2xl border border-white/15 p-6"
            >
              <div className="mb-4 flex items-center gap-3">
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED] text-white">
                  <Sparkles className="size-5" />
                </span>
                <span className="inline-flex h-7 w-7 items-center justify-center rounded-full bg-white/10 text-sm font-semibold text-white">
                  {index + 1}
                </span>
              </div>
              <h3 className="text-lg font-semibold text-white">{step.title}</h3>
              <p className="mt-3 text-sm leading-6 text-[#CBD5E1]">
                {step.description}
              </p>
              <div className="mt-4 flex items-center gap-2 text-xs text-[#CBD5E1]">
                <CheckCircle2 className="size-4 text-emerald-300" />
                Clear and trackable progress
              </div>
            </article>
          ))}
        </div>
      </Container>
    </section>
  );
}

