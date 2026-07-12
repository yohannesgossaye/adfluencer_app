import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { Container } from "@/components/common/container";
import { GradientButton } from "@/components/common/gradient-button";

export function UserTypeCtaSection() {
  return (
    <section id="businesses" className="bg-[color:var(--adfluence-surface)]/50 py-20">
      <Container>
        <div className="grid gap-8 md:grid-cols-2">
          <article className="adfluence-glass relative overflow-hidden rounded-2xl p-10">
            <div className="absolute top-0 right-0 h-64 w-64 rounded-full bg-[color:var(--adfluence-orange)]/15 blur-[100px]" />
            <div className="relative z-10">
              <p className="text-sm text-[#CBD5E1]">For Businesses</p>
              <h3 className="mt-2 mb-4 text-3xl font-bold">For Businesses</h3>
              <p className="mb-8 text-[#CBD5E1]">
                Launch your first campaign and connect with Ethiopia&apos;s top
                creators. Scale your brand with authentic partnerships.
              </p>
              <GradientButton asChild className="h-12 w-full rounded-xl px-6 text-base font-medium sm:text-lg">
                <Link href="/for-businesses">
                  Start Your Campaign <ArrowRight className="size-5" />
                </Link>
              </GradientButton>
            </div>
          </article>

          <article className="adfluence-glass relative overflow-hidden rounded-2xl p-10">
            <div className="absolute top-0 right-0 h-64 w-64 rounded-full bg-[color:var(--adfluence-purple-mid)]/20 blur-[100px]" />
            <div className="relative z-10">
              <p className="text-sm text-[#CBD5E1]">For Influencers</p>
              <h3 className="mt-2 mb-4 text-3xl font-bold">For Influencers</h3>
              <p className="mb-8 text-[#CBD5E1]">
                Join Ethiopia&apos;s premier creator network. Get discovered by
                top brands and monetize your influence with premium campaigns.
              </p>
              <GradientButton asChild className="h-12 w-full rounded-xl px-6 text-base font-medium sm:text-lg">
                <Link href="/for-influencers">
                  Join as Creator <ArrowRight className="size-5" />
                </Link>
              </GradientButton>
            </div>
          </article>
        </div>
      </Container>
    </section>
  );
}
