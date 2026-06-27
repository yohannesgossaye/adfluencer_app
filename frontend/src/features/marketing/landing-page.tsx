import { MarketingFooter } from "@/components/marketing/marketing-footer";
import { FeatureHighlightsSection } from "@/components/marketing/sections/feature-highlights-section";
import { FinalCtaSection } from "@/components/marketing/sections/final-cta-section";
import { HeroSection } from "@/components/marketing/sections/hero-section";
import { HowItWorksSection } from "@/components/marketing/sections/how-it-works-section";
import { StatsSection } from "@/components/marketing/sections/stats-section";
import { TrustedLogosSection } from "@/components/marketing/sections/trusted-logos-section";
import { UserTypeCtaSection } from "@/components/marketing/sections/user-type-cta-section";
import { WhyChooseSection } from "@/components/marketing/sections/why-choose-section";

export function LandingPage() {
  return (
    <div className="min-h-screen bg-[color:var(--adfluence-navy)] text-white">
      <HeroSection />
      <TrustedLogosSection />
      <StatsSection />
      <HowItWorksSection />
      <WhyChooseSection />
      <FeatureHighlightsSection />
      <UserTypeCtaSection />
      <FinalCtaSection />
      <MarketingFooter />
    </div>
  );
}
