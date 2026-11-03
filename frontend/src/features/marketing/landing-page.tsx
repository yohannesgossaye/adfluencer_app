import { HeroDeal } from "@/components/marketing/home/hero-deal";
import { ContactCta, HonestFees, HowADealWorks, TwoSides } from "@/components/marketing/home/home-sections";

/** Home page, Direction 3: Direction 1 "the deal" with one hero photo (option B). */
export function LandingPage() {
  return (
    <div className="bg-[color:var(--adfluence-bg)] text-white">
      <HeroDeal />
      <HowADealWorks />
      <TwoSides />
      <HonestFees />
      <ContactCta />
    </div>
  );
}
