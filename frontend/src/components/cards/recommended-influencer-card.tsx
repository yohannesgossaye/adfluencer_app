import Image from "next/image";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import type { Influencer } from "@/types";

export function RecommendedInfluencerCard({
  influencer,
}: {
  influencer: Influencer;
}) {
  const platform = influencer.platforms[0] ?? "Instagram";

  return (
    <article className="adfluence-glass overflow-hidden rounded-2xl border border-white/15">
      <div className="relative h-44">
        <Image
          src={influencer.imageUrl ?? "/images/adfluence-hero.jpg"}
          alt={influencer.name}
          fill
          className="object-cover"
        />
      </div>
      <div className="space-y-3 p-4 text-white">
        <div className="flex items-center justify-between gap-2">
          <h3 className="text-base font-semibold">{influencer.name}</h3>
          <span className="rounded-full bg-gradient-to-r from-[#FF7A00] to-[#7C3AED] px-3 py-1 text-xs font-medium text-white">
            {influencer.matchScore ?? 90}% Match
          </span>
        </div>
        <p className="text-sm text-[#CBD5E1] capitalize">
          {influencer.niche} / {platform}
        </p>
        <p className="text-sm text-[#CBD5E1]">
          {influencer.followers.toLocaleString()} followers
        </p>
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

