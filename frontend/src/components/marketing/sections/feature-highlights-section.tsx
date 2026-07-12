import {
  BarChart3,
  FileText,
  MessageSquare,
  TrendingUp,
  Users,
  Zap,
} from "lucide-react";

import { Container } from "@/components/common/container";

const highlights = [
  {
    title: "AI Recommendation Engine",
    description:
      "Our smart algorithm matches you with creators that fit your brand and audience.",
    icon: Zap,
  },
  {
    title: "Verified Creator Profiles",
    description:
      "Access detailed profiles with audience insights, engagement, and performance signals.",
    icon: Users,
  },
  {
    title: "Campaign Workspace",
    description:
      "Manage briefs, approvals, and deliverables in one place for both sides.",
    icon: BarChart3,
  },
  {
    title: "Built-in Chat",
    description:
      "Communicate seamlessly with businesses and creators through a focused chat experience.",
    icon: MessageSquare,
  },
  {
    title: "Smart Contracts",
    description:
      "Create, track, and manage agreements with clear milestones and transparency.",
    icon: FileText,
  },
  {
    title: "Advanced Analytics",
    description:
      "Track impressions, engagement, conversions, and ROI with real-time reporting.",
    icon: TrendingUp,
  },
];

export function FeatureHighlightsSection() {
  return (
    <section className="py-20">
      <Container>
        <div className="text-center mb-16">
          <h2 className="text-4xl md:text-5xl font-bold mb-4">
            Powerful Features
          </h2>
          <p className="text-xl text-[#CBD5E1]">
            Everything you need to succeed
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          {highlights.map((item) => (
            <article
              key={item.title}
              className="adfluence-glass rounded-2xl p-8 border border-white/20"
            >
              <span className="mb-6 flex h-14 w-14 items-center justify-center rounded-xl bg-[color:var(--adfluence-purple)] text-white">
                <item.icon className="size-5" />
              </span>
              <h3 className="text-xl font-bold mb-3">{item.title}</h3>
              <p className="text-[#CBD5E1]">{item.description}</p>
            </article>
          ))}
        </div>
      </Container>
    </section>
  );
}
