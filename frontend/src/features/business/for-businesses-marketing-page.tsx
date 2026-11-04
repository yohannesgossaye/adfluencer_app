import { ClosingCta, Faq, NumberedPoints, PageHero, SignupButtons } from "@/components/marketing/marketing-ui";
import { mContainer, mTitle } from "@/components/marketing/marketing-styles";
import { cn } from "@/lib/utils";

const WHY = [
  {
    title: "Only checked creators",
    text: "Every creator is reviewed by our team before they can be invited: real profile, real platform links, and prices they set themselves.",
  },
  {
    title: "Matches that explain themselves",
    text: "Each match shows why it fits your campaign, so you choose with reasons instead of guessing from follower counts.",
  },
  {
    title: "One contract, in birr",
    text: "Deliverables, dates and the price in one agreement. The 5% Adfluence fee is shown before anyone signs.",
  },
  {
    title: "Pay for work you approve",
    text: "Your payment is held by Adfluence and only released to the creator when you approve what they delivered.",
  },
];

// The real weights used by the matching engine (backend helper/matching.go).
const SCORE = [
  { label: "Niche", points: 40, text: "The creator works in your campaign's category" },
  { label: "Platform", points: 25, text: "They are active where you want to be seen" },
  { label: "Budget", points: 20, text: "Their price range overlaps your budget" },
  { label: "Audience size", points: 10, text: "Their following fits the creator size you chose" },
  { label: "Verified", points: 5, text: "Their profile has been verified by our team" },
];

const FAQS = [
  {
    q: "Does my business need to be verified?",
    a: "Yes. To protect creators, we check every business before it can publish campaigns, invite creators or sign contracts. You add your details in your business profile, and our team reviews them, usually within one working day.",
  },
  {
    q: "Can I choose creators myself instead of the matches?",
    a: "Yes. Matches are a shortlist to save you time. You can open any creator's profile, message them and invite the ones you like.",
  },
  {
    q: "What if the creator doesn't deliver?",
    a: "The money stays held until you approve the work. If something goes wrong, our team steps in to help both sides reach a fair outcome.",
  },
  {
    q: "What does it cost?",
    a: "Posting campaigns, matching and chatting are free. When a contract is signed, you pay the agreed price plus a 5% fee. There are no plans or subscriptions.",
  },
];

export function ForBusinessesMarketingPage() {
  return (
    <div className="bg-[color:var(--adfluence-bg)] text-white">
      <PageHero
        title="Run creator campaigns without the guesswork."
        lead="Post one brief, get a ranked shortlist of checked Ethiopian creators, and agree everything in one contract in birr. You only pay for work you approve."
      >
        <SignupButtons role="business" primary="Post a campaign" />
      </PageHero>

      <section className={cn(mContainer, "pb-14 sm:pb-24")}>
        <h2 className={mTitle}>Why businesses use Adfluence</h2>
        <NumberedPoints items={WHY} />
      </section>

      <section className={cn(mContainer, "grid items-start gap-10 pb-14 sm:pb-24 lg:grid-cols-[minmax(0,1fr)_minmax(0,520px)] lg:gap-16")}>
        <div>
          <h2 className={mTitle}>How matching works</h2>
          <p className="mt-4 max-w-[480px] text-base leading-relaxed text-[#CBD5E1] sm:text-lg">
            Every approved creator gets a score out of 100 for your campaign. Only creators who score 78 or more are shown,
            so the shortlist is short and relevant.
          </p>
        </div>
        <div className="rounded-[24px] border border-white/10 bg-[color:var(--adfluence-card)] p-6 sm:p-8">
          <ul className="space-y-5">
            {SCORE.map((s) => (
              <li key={s.label}>
                <div className="flex items-baseline justify-between gap-4">
                  <span className="text-[15px] font-semibold text-white">{s.label}</span>
                  <span className="text-sm text-[#94A3B8] tabular-nums">up to {s.points} points</span>
                </div>
                <div className="mt-2 h-2 rounded-full bg-white/10">
                  <div className="h-2 rounded-full bg-[#A78BFA]" style={{ width: `${(s.points / 40) * 100}%` }} />
                </div>
                <p className="mt-1.5 text-[13px] text-[#94A3B8]">{s.text}</p>
              </li>
            ))}
          </ul>
        </div>
      </section>

      <Faq title="Questions from businesses" items={FAQS} />
      <ClosingCta
        title="Post your first campaign"
        lead="It's free to post. You only pay when a contract is signed."
        role="business"
        primary="Post a campaign"
      />
    </div>
  );
}
