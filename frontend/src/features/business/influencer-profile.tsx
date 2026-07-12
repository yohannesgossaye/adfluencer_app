"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { influencers } from "@/data/influencers";

type InfluencerProfileProps = {
  influencerId: string;
};

const invitedStorageKey = "adfluence_invited_influencers";

const biosByNiche: Record<string, string> = {
  beauty:
    "Creates lifestyle and beauty content for young urban audiences in Ethiopia with strong product storytelling and trusted recommendations.",
  fashion:
    "Shares trend-driven fashion content, styling tips, and creator-led brand integrations optimized for local audiences.",
  food: "Builds food-first social content focused on discovery, culture, and authentic culinary reviews.",
  tech: "Produces practical technology explainers, app reviews, and performance-focused creator content.",
  fitness:
    "Publishes active lifestyle and wellness content designed to drive high engagement and action.",
  lifestyle:
    "Creates relatable day-in-the-life content that blends entertainment and premium brand collaborations.",
  travel:
    "Highlights destination and experience content with strong visual storytelling and discovery reach.",
};

const contentStyleByNiche: Record<string, string> = {
  beauty: "Short-form tutorials, product demos, and before/after formats.",
  fashion: "Outfit reels, trend edits, and premium visual styling.",
  food: "Taste tests, venue highlights, and recipe snippets.",
  tech: "Explainers, comparisons, and honest product walkthroughs.",
  fitness: "Routine breakdowns, challenge formats, and coaching tips.",
  lifestyle: "Daily moments, relatable storytelling, and authentic recommendations.",
  travel: "Destination guides, visual itineraries, and local discovery series.",
};

const audienceByCity: Record<string, string> = {
  "Addis Ababa": "Urban Gen Z and millennial audiences across major Ethiopian cities.",
  "Bahir Dar": "Regional audiences with growing creator engagement in lifestyle and commerce.",
  Hawassa: "Emerging digital audiences interested in technology and creator-driven trends.",
};

function parsePriceRange(priceRangeEtb?: string) {
  if (!priceRangeEtb) return { min: 15000, max: 30000, text: "ETB 15K - 30K" };
  const cleaned = priceRangeEtb.replaceAll("ETB", "").replaceAll(",", "");
  const [minRaw, maxRaw] = cleaned.split("-").map((v) => Number(v.trim().replace("K", "000")));
  const min = Number.isFinite(minRaw) ? minRaw : 15000;
  const max = Number.isFinite(maxRaw) ? maxRaw : 30000;
  return { min, max, text: priceRangeEtb };
}

function getCreatorSizeFromFollowers(followers: number) {
  if (followers < 70000) return "Micro creator";
  if (followers < 110000) return "Mid-tier creator";
  return "Large creator";
}

function getMatchScore(value?: number) {
  return Math.min(96, Math.max(78, value ?? 86));
}

function getFitRows(score: number) {
  return [
    { label: "Niche fit", value: Math.min(98, score + 2) },
    { label: "Platform fit", value: Math.max(76, score - 3) },
    { label: "Budget fit", value: Math.max(74, score - 6) },
    { label: "Audience fit", value: Math.min(99, score + 1) },
    { label: "Verified profile", value: 96 },
  ];
}

export function InfluencerProfile({ influencerId }: InfluencerProfileProps) {
  const [inviteSent, setInviteSent] = useState(false);
  const influencer = influencers.find((entry) => entry.id === influencerId);

  useEffect(() => {
    if (!influencer) return;
    const raw = localStorage.getItem(invitedStorageKey);
    if (!raw) return;
    try {
      const invited = JSON.parse(raw) as string[];
      setInviteSent(invited.includes(influencer.id));
    } catch {
      setInviteSent(false);
    }
  }, [influencer]);

  if (!influencer) {
    return (
      <div className="space-y-6 text-white">
        <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center">
          <h1 className="text-2xl font-semibold text-white">Influencer not found</h1>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            The profile you requested is unavailable in this mock dataset.
          </p>
          <Button
            asChild
            className="mt-5 h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-6 text-white hover:bg-[color:var(--adfluence-orange-hover)]"
          >
            <Link href="/business/recommendations">Back to Recommendations</Link>
          </Button>
        </section>
      </div>
    );
  }

  const score = getMatchScore(influencer.matchScore);
  const price = parsePriceRange(influencer.priceRangeEtb);
  const creatorSize = getCreatorSizeFromFollowers(influencer.followers);
  const responseTime = score >= 92 ? "Within 6 hours" : score >= 86 ? "Within 12 hours" : "Within 24 hours";
  const bio = biosByNiche[influencer.niche] ?? biosByNiche.lifestyle;
  const contentStyle = contentStyleByNiche[influencer.niche] ?? contentStyleByNiche.lifestyle;
  const audience = audienceByCity[influencer.city] ?? "Audience concentrated in Ethiopian urban markets.";
  const avgReach = Math.round(influencer.followers * 0.22);
  const avgViews = Math.round(influencer.followers * 0.18);
  const completedCampaigns = Math.max(4, Math.round(influencer.followers / 16000));
  const fitRows = getFitRows(score);

  const contentSamples = [
    {
      id: "sample-1",
      platform: influencer.platforms[0] ?? "Instagram",
      type: influencer.platforms[0] === "YouTube" ? "YouTube Short" : influencer.platforms[0] === "TikTok" ? "TikTok" : "Reel",
      views: `${Math.round(avgViews * 1.2).toLocaleString()} views`,
      engagement: `${(influencer.engagementRate + 0.7).toFixed(1)}% engagement`,
    },
    {
      id: "sample-2",
      platform: "Instagram",
      type: "Story",
      views: `${Math.round(avgViews * 0.8).toLocaleString()} views`,
      engagement: `${(influencer.engagementRate + 0.3).toFixed(1)}% engagement`,
    },
    {
      id: "sample-3",
      platform: "TikTok",
      type: "TikTok",
      views: `${Math.round(avgViews * 1.35).toLocaleString()} views`,
      engagement: `${(influencer.engagementRate + 0.9).toFixed(1)}% engagement`,
    },
    {
      id: "sample-4",
      platform: "YouTube",
      type: "YouTube Short",
      views: `${Math.round(avgViews * 1.05).toLocaleString()} views`,
      engagement: `${(influencer.engagementRate + 0.4).toFixed(1)}% engagement`,
    },
  ];

  const onInvite = () => {
    const raw = localStorage.getItem(invitedStorageKey);
    let invited: string[] = [];
    if (raw) {
      try {
        invited = JSON.parse(raw) as string[];
      } catch {
        invited = [];
      }
    }
    if (!invited.includes(influencer.id)) {
      invited.push(influencer.id);
      localStorage.setItem(invitedStorageKey, JSON.stringify(invited));
    }
    setInviteSent(true);
  };

  const packageCards = useMemo(
    () => [
      {
        title: "Basic",
        items: ["1 short video", "1 story", "Creator posting included"],
        price: `From ETB ${Math.round(price.min * 0.9).toLocaleString()}`,
      },
      {
        title: "Standard",
        items: ["2 videos", "3 stories", "30-day usage rights"],
        price: `From ETB ${Math.round((price.min + price.max) / 2).toLocaleString()}`,
        recommended: true,
      },
      {
        title: "Premium",
        items: ["Campaign bundle", "Multi-platform package", "Priority revisions"],
        price: `From ETB ${Math.round(price.max * 1.2).toLocaleString()}`,
      },
    ],
    [price.max, price.min],
  );

  return (
    <div className="space-y-6 text-white">
      <section className="grid gap-5 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 lg:grid-cols-[320px_minmax(0,1fr)]">
        <div className="relative h-72 overflow-hidden rounded-2xl border border-white/10 md:h-80">
          <Image
            src={influencer.imageUrl || "/images/influencers/influencer-1.jpg"}
            alt={influencer.name}
            fill
            className="object-cover"
          />
          <span className="absolute right-3 top-3 rounded-full bg-[#0B1A2B]/90 px-2.5 py-1 text-xs font-medium text-emerald-300">
            Verified
          </span>
        </div>

        <div className="space-y-4">
          <div className="flex flex-wrap items-start justify-between gap-3">
            <div>
              <h1 className="text-2xl font-semibold text-white">{influencer.name}</h1>
              <p className="mt-1 text-sm capitalize text-[#CBD5E1]">
                {influencer.niche} • {influencer.city}
              </p>
              <p className="mt-1 text-sm text-[#CBD5E1]">
                Platforms: {influencer.platforms.join(", ")}
              </p>
            </div>
            <span className="rounded-full bg-[color:var(--adfluence-orange)]/15 px-3 py-1 text-xs font-semibold text-orange-200">
              Match score {score}%
            </span>
          </div>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs text-[#CBD5E1]">Followers</p>
              <p className="mt-1 text-lg font-semibold text-white">
                {influencer.followers.toLocaleString()}
              </p>
            </div>
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs text-[#CBD5E1]">Engagement rate</p>
              <p className="mt-1 text-lg font-semibold text-white">{influencer.engagementRate}%</p>
            </div>
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs text-[#CBD5E1]">Price range</p>
              <p className="mt-1 text-lg font-semibold text-white">{price.text}</p>
            </div>
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs text-[#CBD5E1]">Response time</p>
              <p className="mt-1 text-lg font-semibold text-white">{responseTime}</p>
            </div>
          </div>

          <div className="grid gap-2 sm:grid-cols-3">
            <Button
              onClick={onInvite}
              className="h-11 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
            >
              {inviteSent ? "Invite Sent" : "Invite to Campaign"}
            </Button>
            <Button asChild className="h-11 rounded-xl bg-white/10 text-white hover:bg-white/15">
              <Link href="/business/messages">Message</Link>
            </Button>
            <Button asChild className="h-11 rounded-xl bg-white/10 text-white hover:bg-white/15">
              <Link href="/business/contracts">Create Contract</Link>
            </Button>
          </div>
        </div>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
        <h2 className="text-lg font-semibold text-white">About Creator</h2>
        <div className="mt-3 grid gap-3 text-sm text-[#CBD5E1] md:grid-cols-2">
          <p>{bio}</p>
          <p>
            <span className="text-white">Content style:</span> {contentStyle}
          </p>
          <p>
            <span className="text-white">Audience:</span> {audience}
          </p>
          <p>
            <span className="text-white">Languages:</span> Amharic, English
          </p>
          <p className="md:col-span-2">
            <span className="text-white">Brand fit notes:</span> Strong fit for {influencer.niche}{" "}
            campaigns targeting Ethiopian urban audiences and mobile-first social platforms.
          </p>
        </div>
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-white">Performance Metrics</h2>
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {[
            { label: "Followers", value: influencer.followers.toLocaleString() },
            { label: "Engagement Rate", value: `${influencer.engagementRate}%` },
            { label: "Avg Reach", value: avgReach.toLocaleString() },
            { label: "Avg Views", value: avgViews.toLocaleString() },
            { label: "Completed Campaigns", value: completedCampaigns.toString() },
            { label: "Match Score", value: `${score}%` },
          ].map((item) => (
            <article
              key={item.label}
              className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
            >
              <p className="text-xs text-[#CBD5E1]">{item.label}</p>
              <p className="mt-2 text-2xl font-bold text-[color:var(--adfluence-orange)]">
                {item.value}
              </p>
            </article>
          ))}
        </div>
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-white">Content Samples</h2>
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {contentSamples.map((sample, index) => (
            <article
              key={sample.id}
              className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-3"
            >
              <div className="relative h-36 overflow-hidden rounded-xl">
                <Image
                  src={influencer.imageUrl || "/images/influencers/influencer-1.jpg"}
                  alt={`${influencer.name} sample ${index + 1}`}
                  fill
                  className="object-cover"
                />
                <span className="absolute left-2 top-2 rounded-full bg-[color:var(--adfluence-surface)]/90 px-2 py-0.5 text-xs text-[#CBD5E1]">
                  {sample.platform}
                </span>
              </div>
              <p className="mt-3 text-sm font-medium text-white">{sample.type}</p>
              <p className="mt-1 text-xs text-[#CBD5E1]">{sample.views}</p>
              <p className="text-xs text-[#CBD5E1]">{sample.engagement}</p>
            </article>
          ))}
        </div>
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-white">Pricing &amp; Packages</h2>
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {packageCards.map((pkg) => (
            <article
              key={pkg.title}
              className={`rounded-2xl border p-5 ${
                pkg.recommended
                  ? "border-[color:var(--adfluence-orange)]/40 bg-[color:var(--adfluence-orange)]/10"
                  : "border-white/10 bg-[rgba(255,255,255,0.08)]"
              }`}
            >
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold text-white">{pkg.title}</h3>
                {pkg.recommended ? (
                  <span className="rounded-full bg-white/15 px-2 py-0.5 text-xs text-white">
                    Recommended
                  </span>
                ) : null}
              </div>
              <ul className="mt-3 space-y-1 text-sm text-[#CBD5E1]">
                {pkg.items.map((item) => (
                  <li key={item}>• {item}</li>
                ))}
              </ul>
              <p className="mt-4 text-lg font-semibold text-white">{pkg.price}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="grid gap-6 xl:grid-cols-2">
        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <h2 className="text-lg font-semibold text-white">Campaign Fit Analysis</h2>
          <div className="mt-4 space-y-3">
            {fitRows.map((row) => (
              <div key={row.label}>
                <div className="mb-1 flex items-center justify-between text-xs">
                  <span className="text-[#CBD5E1]">{row.label}</span>
                  <span className="font-medium text-white">{row.value}%</span>
                </div>
                <div className="h-2 rounded-full bg-white/10">
                  <div
                    className="h-2 rounded-full bg-[color:var(--adfluence-orange)]"
                    style={{ width: `${row.value}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <h2 className="text-lg font-semibold text-white">Reviews / Past Collaborations</h2>
          <div className="mt-4 space-y-3 text-sm">
            <article className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="font-medium text-white">Betu Cosmetics</p>
              <p className="mt-1 text-[#CBD5E1]">Rating: 4.9/5 • Campaign delivered on time.</p>
            </article>
            <article className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="font-medium text-white">Habesha Coffee Roasters</p>
              <p className="mt-1 text-[#CBD5E1]">Rating: 4.8/5 • Strong engagement and reach uplift.</p>
            </article>
            <article className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="font-medium text-white">Ride Ethiopia</p>
              <p className="mt-1 text-[#CBD5E1]">Rating: 4.7/5 • High CTR on creator-led campaign assets.</p>
            </article>
          </div>
        </div>
      </section>
    </div>
  );
}
