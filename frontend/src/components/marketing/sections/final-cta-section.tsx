import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { Container } from "@/components/common/container";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";

export function FinalCtaSection() {
  return (
    <section className="py-24 bg-gradient-to-r from-[#FF7A00] via-[#7C3AED] to-[#FF7A00] relative overflow-hidden">
      <div className="absolute inset-0 bg-[color:var(--adfluence-navy)]/40" />
      <Container className="relative z-10 text-center text-white">
        <div className="mx-auto max-w-4xl">
        <h2 className="text-5xl md:text-6xl font-bold mb-6">
          Ready to take your brand to the next level?
        </h2>
        <p className="text-xl mb-10 opacity-90">
          Let&apos;s create something amazing together.
        </p>
        <div className="flex flex-wrap gap-4 justify-center">
          <Button
            asChild
            className="h-12 rounded-xl bg-white px-6 text-base font-medium text-[color:var(--adfluence-text-dark)] hover:bg-slate-50 sm:text-lg"
          >
            <Link href="/auth/signup">
              Get Started Free <ArrowRight className="size-5" />
            </Link>
          </Button>
          <Button
            asChild
            variant="outline"
            className="h-12 rounded-xl border-white/25 bg-transparent px-6 text-base font-medium text-white hover:bg-white/10 sm:text-lg"
          >
            <Link href="/auth/login">Schedule a Demo</Link>
          </Button>
        </div>
        </div>
      </Container>
    </section>
  );
}
