import {
  AlertTriangle,
  Handshake,
  MessageSquare,
  Search,
  ShieldCheck,
  Timer,
} from "lucide-react";

import { Container } from "@/components/common/container";

const problems = [
  {
    title: "Fake engagement",
    description:
      "Avoid inflated numbers with creator signals and verification workflows.",
    icon: AlertTriangle,
  },
  {
    title: "Manual outreach",
    description:
      "Stop chasing DMs. Discover and shortlist creators in minutes, not days.",
    icon: Search,
  },
  {
    title: "Scattered communication",
    description:
      "Keep feedback, approvals, and briefs in one messaging workspace.",
    icon: MessageSquare,
  },
  {
    title: "Poor tracking",
    description:
      "Track performance and ROI with clear reporting and campaign status.",
    icon: Timer,
  },
  {
    title: "Contract & payment risk",
    description:
      "Use structured agreements and secure collaboration for both parties.",
    icon: ShieldCheck,
  },
  {
    title: "Hard creator discovery",
    description:
      "Match the right niche, platform, and audience for each campaign goal.",
    icon: Handshake,
  },
];

export function ProblemsWeSolveSection() {
  return (
    <section className="py-20 bg-[color:var(--adfluence-navy-dark)]/50">
      <Container>
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-white sm:text-4xl">
            Problems we solve for modern teams
          </h2>
          <p className="mt-3 text-[#CBD5E1]">
            Practical fixes for the toughest parts of influencer marketing.
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {problems.map((p) => (
            <article
              key={p.title}
              className="adfluence-glass rounded-2xl border border-white/15 p-6 transition-transform hover:-translate-y-0.5"
            >
              <div className="flex items-center gap-3">
                <span className="inline-flex h-11 w-11 items-center justify-center rounded-xl bg-[color:var(--adfluence-purple)] text-white">
                  <p.icon className="size-5" />
                </span>
                <h3 className="text-base font-semibold text-white">{p.title}</h3>
              </div>
              <p className="mt-3 text-sm leading-6 text-[#CBD5E1]">
                {p.description}
              </p>
            </article>
          ))}
        </div>
      </Container>
    </section>
  );
}

