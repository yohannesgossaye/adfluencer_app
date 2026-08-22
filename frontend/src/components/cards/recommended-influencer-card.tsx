import Image from "next/image";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { getInitials } from "@/lib/utils";
import type { Influencer } from "@/types";

export type RecommendedInfluencerCardData = {
  id: string;
  name: string;
  imageUrl?: string;
  niche?: string;
  platform?: string;
  followers?: number;
  priceRangeEtb?: string;
  matchScore?: number;
};

export function RecommendedInfluencerCard({
  influencer,
}: {
  influencer: RecommendedInfluencerCardData | Influencer;
}) {
  const platform =
    "platform" in influencer && influencer.platform
      ? influencer.platform
      : "platforms" in influencer && influencer.platforms?.[0]
        ? influencer.platforms[0]
        : "Social Media";

  const followersCount =
    typeof influencer.followers === "number" ? influencer.followers.toLocaleString() : "Active";

  return (
    <article className="adfluence-glass overflow-hidden rounded-2xl border border-white/15">
      <div className="relative flex h-44 items-center justify-center overflow-hidden bg-gradient-to-br from-[#0B1A2B] via-white/5 to-[color:var(--adfluence-orange)]/20">
        {influencer.imageUrl ? (
          <Image
            src={influencer.imageUrl}
            alt={influencer.name}
            fill
            className="object-cover"
          />
        ) : (
          <div className="flex flex-col items-center justify-center gap-1">
            <span className="text-3xl font-extrabold text-white uppercase tracking-wider">
              {getInitials(influencer.name, "CR")}
            </span>
            <span className="text-xs text-[#94A3B8]">
              {influencer.niche || "Creator"}
            </span>
          </div>
        )}
      </div>
      <div className="space-y-3 p-4 text-white">
        <div className="flex items-center justify-between gap-2">
          <h3 className="text-base font-semibold truncate">{influencer.name}</h3>
          <span className="shrink-0 rounded-full bg-[color:var(--adfluence-orange)] px-3 py-1 text-xs font-medium text-white">
            {influencer.matchScore ?? 90}% Match
          </span>
        </div>
        <p className="text-sm text-[#CBD5E1] capitalize truncate">
          {influencer.niche || "Creator"} / {platform}
        </p>
        <p className="text-sm text-[#CBD5E1]">{followersCount} followers</p>
        <p className="text-sm text-[#CBD5E1]">
          {influencer.priceRangeEtb ?? "ETB 15K - 30K"}
        </p>
        <Button
          asChild
          className="h-10 w-full rounded-xl bg-white/10 text-white hover:bg-white/15"
        >
          <Link href={`/business/influencers/${influencer.id}`}>View Profile</Link>
        </Button>
      </div>
    </article>
  );
}
