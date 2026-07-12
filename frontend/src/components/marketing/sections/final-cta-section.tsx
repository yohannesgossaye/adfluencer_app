import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { Container } from "@/components/common/container";
import { Button } from "@/components/ui/button";
import { buildMarketingAuthHref } from "@/lib/marketing-auth-url";

export function FinalCtaSection() {
  return (
    <section className="relative overflow-hidden bg-[color:var(--adfluence-surface)] py-24">
      <div className="absolute top-1/2 left-1/2 h-96 w-96 -translate-x-1/2 -translate-y-1/2 rounded-full bg-[color:var(--adfluence-orange)]/8 blur-[120px]" />
      <Container className="relative z-10 text-center text-white">
        <div className="mx-auto max-w-4xl">
          <h2 className="mb-6 text-5xl font-bold md:text-6xl">
            Ready to take your brand to the next level?
          </h2>
          <p className="mb-10 text-xl opacity-90">
            Let&apos;s create something amazing together.
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            <Button
              asChild
              className="h-12 rounded-xl bg-[color:var(--adfluence-orange)] px-6 text-base font-medium text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:text-lg"
            >
              <Link href={buildMarketingAuthHref("/", "signup")}>
                Get Started Free <ArrowRight className="size-5" />
              </Link>
            </Button>
            <Button
              asChild
              variant="outline"
              className="h-12 rounded-xl border-white/25 bg-transparent px-6 text-base font-medium text-white hover:bg-white/10 sm:text-lg"
            >
              <Link href={buildMarketingAuthHref("/", "login")}>
                Schedule a Demo
              </Link>
            </Button>
          </div>
        </div>
      </Container>
    </section>
  );
}
