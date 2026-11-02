import Image from "next/image";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { getInitials } from "@/lib/utils";

export type RecommendedInfluencerCardData = {
  id: string;
  name: string;
  imageUrl?: string;
  matchScore: number;
  campaignId?: string;
  campaignTitle?: string;
};

/** A creator recommended for one of the business's campaigns. Shows only data the API returns. */
export function RecommendedInfluencerCard({ influencer }: { influencer: RecommendedInfluencerCardData }) {
  const href = influencer.campaignId
    ? `/business/influencers/${influencer.id}?campaign_id=${influencer.campaignId}`
    : `/business/influencers/${influencer.id}`;

  return (
    <article className="adfluence-glass overflow-hidden rounded-2xl border border-white/15">
      <div className="relative flex h-44 items-center justify-center overflow-hidden bg-gradient-to-br from-[#0B1A2B] via-white/5 to-[color:var(--adfluence-orange)]/20">
        {influencer.imageUrl ? (
          <Image src={influencer.imageUrl} alt={influencer.name} fill className="object-cover" />
        ) : (
          <span className="text-3xl font-extrabold uppercase tracking-wider text-white" aria-hidden>
            {getInitials(influencer.name, "CR")}
          </span>
        )}
      </div>
      <div className="space-y-3 p-4 text-white">
        <div className="flex items-center justify-between gap-2">
          <h3 className="truncate text-base font-semibold">{influencer.name}</h3>
          <span className="shrink-0 rounded-full bg-[color:var(--adfluence-orange)] px-3 py-1 text-xs font-medium text-white">
            {Math.round(influencer.matchScore)}% match
          </span>
        </div>
        {influencer.campaignTitle ? (
          <p className="truncate text-sm text-[#CBD5E1]">For: {influencer.campaignTitle}</p>
        ) : null}
        <Button asChild className="h-10 w-full rounded-xl bg-white/10 text-white hover:bg-white/15">
          <Link href={href}>View profile</Link>
        </Button>
      </div>
    </article>
  );
}
