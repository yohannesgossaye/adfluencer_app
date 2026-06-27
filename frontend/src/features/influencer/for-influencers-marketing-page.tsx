import {
  InfluencerFeatureHighlightsSection,
  InfluencerFinalCtaSection,
  InfluencerHeroSection,
  InfluencerHowItWorksSection,
  InfluencerTrustStatsSection,
  InfluencerWhyLoveSection,
} from "@/components/marketing/sections";

export function ForInfluencersMarketingPage() {
  return (
    <>
      <InfluencerHeroSection />
      <InfluencerTrustStatsSection />
      <InfluencerHowItWorksSection />
      <InfluencerWhyLoveSection />
      <InfluencerFeatureHighlightsSection />
      <InfluencerFinalCtaSection />
    </>
  );
}

