import {
  BarChart3,
  CreditCard,
  FileText,
  LayoutTemplate,
  MessageSquare,
  Search,
  ShieldCheck,
  Users,
} from "lucide-react";

import { Container } from "@/components/common/container";

const features = [
  { title: "AI Creator Matching", icon: Search },
  { title: "Campaign Builder", icon: LayoutTemplate },
  { title: "Messaging Workspace", icon: MessageSquare },
  { title: "Performance Analytics", icon: BarChart3 },
  { title: "Contract Management", icon: FileText },
  { title: "Secure Payments", icon: CreditCard },
  { title: "ROI Reporting", icon: ShieldCheck },
  { title: "Audience Insights", icon: Users },
];

export function PlatformFeaturesSection() {
  return (
    <section className="py-20">
      <Container>
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-white sm:text-4xl">
            Business platform features
          </h2>
          <p className="mt-3 text-[#CBD5E1]">
            Everything you need to launch, manage, and scale creator campaigns.
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {features.map((f) => (
            <article
              key={f.title}
              className="adfluence-glass rounded-2xl border border-white/15 p-6"
            >
              <span className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-[color:var(--adfluence-purple)] text-white">
                <f.icon className="size-5" />
              </span>
              <h3 className="mt-4 text-base font-semibold text-white">{f.title}</h3>
              <p className="mt-2 text-sm text-[#CBD5E1]">
                Built for fast execution with clear reporting.
              </p>
            </article>
          ))}
        </div>
      </Container>
    </section>
  );
}

