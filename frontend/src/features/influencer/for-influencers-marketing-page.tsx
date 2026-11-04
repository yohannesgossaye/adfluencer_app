import { ClosingCta, Faq, NumberedPoints, PageHero, SignupButtons } from "@/components/marketing/marketing-ui";
import { mContainer, mTitle } from "@/components/marketing/marketing-styles";
import { formatEtb, PLATFORM_FEE_PERCENT, previewFees } from "@/lib/contracts";
import { cn } from "@/lib/utils";

const STEPS = [
  { title: "Create your profile", text: "Add your niche, your TikTok, Instagram or YouTube links, your audience and your prices. Joining is free." },
  { title: "Get approved", text: "Our team reviews your profile so businesses know every creator here is real. If something is missing, we tell you exactly what to fix." },
  { title: "Receive real briefs", text: "Businesses invite you to campaigns that match your niche, platform and prices. Accept the ones you like, decline the rest." },
  { title: "Agree in one contract", text: "Deliverables, dates and the price are written down before you start. The business's payment is secured first." },
  { title: "Deliver and get paid", text: "When the business approves your work, the payment is released to you, minus the 5% Adfluence fee." },
];

// Same arithmetic as contracts (lib/contracts previewFees).
const EXAMPLE = previewFees(20000);

const FAQS = [
  {
    q: "Who can join as a creator?",
    a: "Anyone in Ethiopia creating content on TikTok, Instagram or YouTube. We review each profile, its links and its content before it is shown to businesses.",
  },
  {
    q: "Do I pay anything to join?",
    a: "No. Joining, receiving invites and chatting are free. The 5% fee only comes out of the payment for a contract you finished.",
  },
  {
    q: "Can I say no to an invite?",
    a: "Yes. You choose which campaigns to accept. Declining is free and the business is simply told you're not available.",
  },
  {
    q: "How do I know I'll be paid?",
    a: "Work only starts after the business's payment is held by Adfluence. When the business approves your work, that money is released to you.",
  },
  {
    q: "Can I change my prices later?",
    a: "Yes, any time in your profile. New invites use your updated prices; contracts you already signed keep the agreed price.",
  },
];

export function ForInfluencersMarketingPage() {
  return (
    <div className="bg-[color:var(--adfluence-bg)] text-white">
      <PageHero
        title="Get paid for the content you already make."
        lead="Join free, get approved by our team, and receive real campaign briefs from Ethiopian businesses, with the price agreed and the payment secured before you start."
      >
        <SignupButtons role="influencer" primary="Create my creator profile" />
      </PageHero>

      <section className={cn(mContainer, "pb-14 sm:pb-24")}>
        <h2 className={mTitle}>How it works for creators</h2>
        <NumberedPoints items={STEPS} />
      </section>

      <section className={cn(mContainer, "grid items-center gap-10 pb-14 sm:pb-24 lg:grid-cols-[minmax(0,1fr)_minmax(0,480px)] lg:gap-16")}>
        <div>
          <h2 className={mTitle}>Your money, clearly</h2>
          <p className="mt-4 max-w-[480px] text-base leading-relaxed text-[#CBD5E1] sm:text-lg">
            You see the price before you say yes. The business pays it into Adfluence first, and you receive it when they
            approve the work, minus {PLATFORM_FEE_PERCENT}%. No surprises and no chasing payments.
          </p>
        </div>
        <div className="rounded-[24px] border border-white/10 bg-[color:var(--adfluence-card)] p-6 sm:p-8">
          <p className="text-sm font-semibold text-[#A78BFA]">Example: a {formatEtb(EXAMPLE.amount)} contract</p>
          <dl className="mt-5 space-y-4">
            <div className="flex items-baseline justify-between gap-4">
              <dt className="text-[15px] text-[#CBD5E1]">Agreed price</dt>
              <dd className="text-lg font-semibold text-white tabular-nums">{formatEtb(EXAMPLE.amount)}</dd>
            </div>
            <div className="flex items-baseline justify-between gap-4">
              <dt className="text-[15px] text-[#CBD5E1]">Adfluence fee ({PLATFORM_FEE_PERCENT}%)</dt>
              <dd className="text-lg font-semibold text-white tabular-nums">− {formatEtb(EXAMPLE.influencerFee)}</dd>
            </div>
            <div className="flex items-baseline justify-between gap-4 border-t border-white/10 pt-4">
              <dt className="text-[15px] font-semibold text-white">You receive</dt>
              <dd className="text-2xl font-bold text-[color:var(--adfluence-orange)] tabular-nums">{formatEtb(EXAMPLE.influencerReceives)}</dd>
            </div>
          </dl>
          <p className="mt-4 text-[13px] text-[#94A3B8]">Released when the business approves your work.</p>
        </div>
      </section>

      <Faq title="Questions from creators" items={FAQS} />
      <ClosingCta
        title="Start your creator profile"
        lead="Free to join. You only pay 5% on deals you finish."
        role="influencer"
        primary="Create my creator profile"
      />
    </div>
  );
}
