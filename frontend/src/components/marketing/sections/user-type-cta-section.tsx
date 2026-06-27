import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { Container } from "@/components/common/container";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";

export function UserTypeCtaSection() {
  return (
    <section id="businesses" className="py-20 bg-[color:var(--adfluence-navy-dark)]/50">
      <Container>
        <div className="grid md:grid-cols-2 gap-8">
          <article className="adfluence-glass p-10 relative overflow-hidden rounded-2xl">
            <div className="absolute top-0 right-0 w-64 h-64 bg-[#FF7A00]/20 rounded-full blur-[100px]" />
            <div className="relative z-10">
              <p className="text-sm text-[#CBD5E1]">For Businesses</p>
              <h3 className="text-3xl font-bold mt-2 mb-4">For Businesses</h3>
              <p className="text-[#CBD5E1] mb-8">
                Launch your first campaign and connect with Ethiopia&apos;s top
                creators. Scale your brand with authentic partnerships.
              </p>
              <GradientButton asChild className="w-full h-12 rounded-xl px-6 text-base font-medium sm:text-lg">
                <Link href="/for-businesses">
                  Start Your Campaign <ArrowRight className="size-5" />
                </Link>
              </GradientButton>
            </div>
          </article>

          <article className="adfluence-glass p-10 relative overflow-hidden rounded-2xl">
            <div className="absolute top-0 right-0 w-64 h-64 bg-[#7C3AED]/20 rounded-full blur-[100px]" />
            <div className="relative z-10">
              <p className="text-sm text-[#CBD5E1]">For Influencers</p>
              <h3 className="text-3xl font-bold mt-2 mb-4">For Influencers</h3>
              <p className="text-[#CBD5E1] mb-8">
                Join Ethiopia&apos;s premier creator network. Get discovered by
                top brands and monetize your influence with premium campaigns.
              </p>
              <Button
                asChild
                className="w-full h-12 rounded-xl bg-gradient-to-r from-[#FF7A00] to-[#7C3AED] px-6 text-base font-medium text-white hover:opacity-95 sm:text-lg"
              >
                <Link href="/for-influencers">
                  Join as Creator <ArrowRight className="size-5" />
                </Link>
              </Button>
            </div>
          </article>
        </div>
      </Container>
    </section>
  );
}
