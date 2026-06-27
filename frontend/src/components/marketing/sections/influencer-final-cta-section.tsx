import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { Container } from "@/components/common/container";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";

export function InfluencerFinalCtaSection() {
  return (
    <section className="relative overflow-hidden bg-gradient-to-r from-[#FF7A00] via-[#7C3AED] to-[#FF7A00] py-24">
      <div className="absolute inset-0 bg-[#0B1A2B]/45" />
      <Container className="relative text-center text-white">
        <div className="mx-auto max-w-3xl">
          <h2 className="text-4xl font-bold leading-tight sm:text-5xl">
            Start growing your creator career today.
          </h2>
          <p className="mt-4 text-lg text-white/90">
            Join campaigns from verified brands and build predictable income with
            Adfluence.
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <GradientButton asChild className="h-12 rounded-xl px-6">
              <Link href="/auth/signup?role=influencer">
                Join as Influencer <ArrowRight className="size-5" />
              </Link>
            </GradientButton>
            <Button
              asChild
              variant="outline"
              className="h-12 rounded-xl border-white/25 bg-transparent px-6 text-white hover:bg-white/10"
            >
              <Link href="/business/campaigns">Explore Campaigns</Link>
            </Button>
          </div>
        </div>
      </Container>
    </section>
  );
}

