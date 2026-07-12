import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { Container } from "@/components/common/container";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { buildMarketingAuthHref } from "@/lib/marketing-auth-url";

export function BusinessFinalCtaSection() {
  return (
    <section className="relative overflow-hidden bg-[color:var(--adfluence-surface)] py-24">
      <div className="absolute top-1/2 left-1/2 h-96 w-96 -translate-x-1/2 -translate-y-1/2 rounded-full bg-[color:var(--adfluence-orange)]/8 blur-[120px]" />
      <Container className="relative text-center text-white">
        <div className="mx-auto max-w-3xl">
          <h2 className="text-4xl font-bold leading-tight sm:text-5xl">
            Ready to scale your influencer marketing?
          </h2>
          <p className="mt-4 text-lg text-white/90">
            Launch your next campaign with verified creators and measurable ROI.
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <GradientButton asChild className="h-12 rounded-xl px-6">
              <Link href={buildMarketingAuthHref("/for-businesses", "signup", "business")}>
                Start Campaign <ArrowRight className="size-5" />
              </Link>
            </GradientButton>
            <Button
              asChild
              variant="outline"
              className="h-12 rounded-xl border-white/25 bg-transparent px-6 text-white hover:bg-white/10"
            >
              <Link href={buildMarketingAuthHref("/for-businesses", "login")}>
                Book Demo
              </Link>
            </Button>
          </div>
        </div>
      </Container>
    </section>
  );
}
