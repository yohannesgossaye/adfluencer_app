import {
  BusinessDashboardPreviewSection,
  BusinessFinalCtaSection,
  BusinessHeroSection,
  CampaignWorkflowSection,
  PlatformFeaturesSection,
  ProblemsWeSolveSection,
  TrustResultsSection,
} from "@/components/marketing/sections";

export function ForBusinessesMarketingPage() {
  return (
    <>
      <BusinessHeroSection />
      <ProblemsWeSolveSection />
      <PlatformFeaturesSection />
      <CampaignWorkflowSection />
      <BusinessDashboardPreviewSection />
      <TrustResultsSection />
      <BusinessFinalCtaSection />
    </>
  );
}

