import Link from "next/link";

import { Container } from "@/components/common/container";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const plans = [
  {
    name: "Starter",
    price: "Free",
    description: "For small businesses testing influencer campaigns.",
    features: [
      "Create 1 campaign",
      "Browse verified influencers",
      "Basic recommendations",
      "Email support",
    ],
    cta: "Start Free",
  },
  {
    name: "Growth",
    price: "999 ETB / month",
    description: "For active brands running monthly campaigns.",
    features: [
      "Unlimited campaigns",
      "Smart influencer recommendations",
      "In-app chat",
      "Contracts",
      "Campaign analytics",
      "Priority support",
    ],
    cta: "Choose Growth",
    popular: true,
  },
  {
    name: "Agency",
    price: "Custom",
    description: "For agencies and larger businesses managing many campaigns.",
    features: [
      "Multi-brand workspace",
      "Advanced analytics",
      "Dedicated support",
      "Custom onboarding",
      "Team access",
      "Payment workflow support",
    ],
    cta: "Contact Sales",
  },
];

const faqs = [
  {
    question: "Can I start free?",
    answer: "Yes. The Starter plan is free for testing your first campaign.",
  },
  {
    question: "Do influencers pay to join?",
    answer: "No. Influencers can join for free and earn via collaborations.",
  },
  {
    question: "Does Adfluence handle payments?",
    answer: "Yes, we support secure payment workflows for collaborations.",
  },
  {
    question: "Can I upgrade later?",
    answer: "Absolutely. You can move from Starter to Growth anytime.",
  },
  {
    question: "Is this built for Ethiopia?",
    answer: "Yes. Adfluence is tailored for Ethiopia's creator and brand market.",
  },
];

export function PricingPage() {
  return (
    <div className="bg-[color:var(--adfluence-navy)] text-white">
      <section className="py-20">
        <Container>
          <div className="mx-auto max-w-3xl text-center">
            <h1 className="text-4xl font-bold sm:text-5xl">
              Simple pricing for growing campaigns.
            </h1>
            <p className="mt-4 text-lg text-[#CBD5E1]">
              Choose the plan that fits your brand, creator journey, or platform
              needs.
            </p>
          </div>
        </Container>
      </section>

      <section className="pb-20">
        <Container>
          <div className="grid gap-5 lg:grid-cols-3">
            {plans.map((plan) => (
              <Card
                key={plan.name}
                className={`relative rounded-2xl border bg-[#0B1A2B]/75 text-white backdrop-blur-xl ${
                  plan.popular
                    ? "border-[#7C3AED] shadow-[0_24px_60px_rgba(124,58,237,0.25)]"
                    : "border-white/15"
                }`}
              >
                {plan.popular ? (
                  <span className="absolute right-4 top-4 rounded-full bg-gradient-to-r from-[#FF7A00] to-[#7C3AED] px-3 py-1 text-xs font-medium text-white">
                    Most Popular
                  </span>
                ) : null}
                <CardHeader className="space-y-3">
                  <CardTitle className="text-2xl text-white">{plan.name}</CardTitle>
                  <p className="text-sm text-[#CBD5E1]">{plan.description}</p>
                  <p className="text-2xl font-semibold text-white">{plan.price}</p>
                </CardHeader>
                <CardContent className="space-y-3">
                  <ul className="space-y-2 text-sm text-[#CBD5E1]">
                    {plan.features.map((feature) => (
                      <li key={feature}>• {feature}</li>
                    ))}
                  </ul>
                  {plan.popular ? (
                    <GradientButton className="mt-4 h-12 w-full rounded-xl px-6">
                      {plan.cta}
                    </GradientButton>
                  ) : (
                    <Button
                      variant="outline"
                      className="mt-4 h-12 w-full rounded-xl border-white/15 bg-transparent px-6 text-white hover:bg-white/10"
                    >
                      {plan.cta}
                    </Button>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </Container>
      </section>

      <section className="pb-16">
        <Container>
          <div className="rounded-2xl border border-white/15 bg-[#0B1A2B]/70 p-6 text-center text-[#CBD5E1] backdrop-blur-xl">
            Influencers can join for free and earn from brand collaborations.
          </div>
        </Container>
      </section>

      <section className="pb-20">
        <Container>
          <div className="mx-auto max-w-4xl">
            <h2 className="text-center text-3xl font-bold text-white">FAQ</h2>
            <div className="mt-8 space-y-3">
              {faqs.map((faq) => (
                <div
                  key={faq.question}
                  className="rounded-2xl border border-white/15 bg-[#0B1A2B]/70 p-5 backdrop-blur-xl"
                >
                  <h3 className="text-lg font-semibold text-white">{faq.question}</h3>
                  <p className="mt-2 text-sm text-[#CBD5E1]">{faq.answer}</p>
                </div>
              ))}
            </div>
          </div>
        </Container>
      </section>

      <section className="relative overflow-hidden bg-gradient-to-r from-[#FF7A00] via-[#7C3AED] to-[#FF7A00] py-24">
        <div className="absolute inset-0 bg-[#0B1A2B]/45" />
        <Container className="relative text-center text-white">
          <h2 className="text-4xl font-bold sm:text-5xl">
            Ready to launch your first campaign?
          </h2>
          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <GradientButton asChild className="h-12 rounded-xl px-6">
              <Link href="/auth/signup">Start Free</Link>
            </GradientButton>
            <Button
              asChild
              variant="outline"
              className="h-12 rounded-xl border-white/25 bg-transparent px-6 text-white hover:bg-white/10"
            >
              <Link href="/auth/login">Talk to Us</Link>
            </Button>
          </div>
        </Container>
      </section>
    </div>
  );
}
