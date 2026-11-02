import { HeroDeal } from "@/components/marketing/home/hero-deal";
import { HonestFees, HowADealWorks, TalkToUs, TwoSides } from "@/components/marketing/home/home-sections";
import { TibebBand } from "@/components/marketing/tibeb-band";

/** Home page, Direction 1 "the deal" (Figma: brand & home page proposal, page 2). */
export function LandingPage() {
  return (
    <div className="bg-[color:var(--adfluence-bg)] text-white">
      <HeroDeal />
      <TibebBand />
      <HowADealWorks />
      <TwoSides />
      <HonestFees />
      <TalkToUs />
    </div>
  );
}
