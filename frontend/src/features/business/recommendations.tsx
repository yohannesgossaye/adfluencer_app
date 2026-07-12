"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { influencers } from "@/data/influencers";

type CampaignDraft = {
  title?: string;
  goal?: string;
  niche?: string;
  platforms?: string[];
  budgetRange?: string;
  creatorSize?: string;
  startDate?: string;
  endDate?: string;
};

type RecommendedInfluencer = (typeof influencers)[number] & {
  calculatedScore: number;
  bio: string;
  estimatedPriceMin: number;
  estimatedPriceMax: number;
  creatorSize: "Micro creator" | "Mid-tier creator" | "Large creator";
};

const invitedStorageKey = "adfluence_invited_influencers";
const draftStorageKey = "adfluence_campaign_draft";

const nicheOptions = [
  "All niches",
  "Beauty",
  "Fashion",
  "Food",
  "Tech",
  "Fitness",
  "Lifestyle",
  "Travel",
] as const;

const platformOptions = ["All platforms", "TikTok", "Instagram", "YouTube"] as const;
const scoreOptions = ["Any score", "80+", "85+", "90+", "95+"] as const;
const sortOptions = [
  "Best match",
  "Followers",
  "Price low to high",
  "Price high to low",
] as const;

const influencerBios: Record<string, string> = {
  "inf-1": "Beauty creator with strong short-form storytelling and product demo performance.",
  "inf-2": "Fashion and lifestyle creator known for high-retention audience engagement.",
  "inf-3": "Style and trend creator with premium content quality and strong conversion potential.",
  "inf-4": "Lifestyle creator with authentic local audience trust and consistent campaign delivery.",
  "inf-5": "Tech creator focused on practical reviews and app-focused creator campaigns.",
};

function getCreatorSizeFromFollowers(followers: number) {
  if (followers < 70000) return "Micro creator";
  if (followers < 110000) return "Mid-tier creator";
  return "Large creator";
}

function getBudgetBand(text?: string) {
  if (!text) return { min: 10000, max: 30000 };
  if (text.includes("Under 10,000")) return { min: 0, max: 10000 };
  if (text.includes("10,000")) return { min: 10000, max: 25000 };
  if (text.includes("25,000")) return { min: 25000, max: 50000 };
  return { min: 50000, max: 150000 };
}

function parsePriceRange(priceRangeEtb?: string) {
  if (!priceRangeEtb) return { min: 15000, max: 30000 };
  const cleaned = priceRangeEtb.replaceAll("ETB", "").replaceAll(",", "");
  const [minRaw, maxRaw] = cleaned.split("-").map((v) => Number(v.trim().replace("K", "000")));
  return {
    min: Number.isFinite(minRaw) ? minRaw : 15000,
    max: Number.isFinite(maxRaw) ? maxRaw : 30000,
  };
}

function formatTimeline(startDate?: string, endDate?: string) {
  if (!startDate || !endDate) return "Not set";
  return `${new Date(startDate).toLocaleDateString()} - ${new Date(endDate).toLocaleDateString()}`;
}

function scoreInfluencer(
  influencer: (typeof influencers)[number],
  draft: CampaignDraft | null,
): RecommendedInfluencer {
  const price = parsePriceRange(influencer.priceRangeEtb);
  const creatorSize = getCreatorSizeFromFollowers(influencer.followers);

  if (!draft) {
    return {
      ...influencer,
      calculatedScore: influencer.matchScore ?? 84,
      bio: influencerBios[influencer.id] ?? "Adfluence creator profile with strong campaign fit.",
      estimatedPriceMin: price.min,
      estimatedPriceMax: price.max,
      creatorSize,
    };
  }

  let score = 0;
  const draftNiche = draft.niche?.toLowerCase();
  const nicheMatch = draftNiche && influencer.niche.toLowerCase() === draftNiche;
  if (nicheMatch) score += 40;

  const platformMatch =
    draft.platforms?.some((platform) => influencer.platforms.includes(platform as "TikTok")) ??
    false;
  if (platformMatch) score += 25;

  const budgetBand = getBudgetBand(draft.budgetRange);
  const withinBudget = price.min <= budgetBand.max && price.max >= budgetBand.min;
  if (withinBudget) score += 20;

  const creatorSizeFit = draft.creatorSize
    ? creatorSize.toLowerCase() === draft.creatorSize.toLowerCase()
    : false;
  if (creatorSizeFit) score += 10;

  if (influencer.status === "approved") score += 5;

  if (score < 78) {
    const fallback = influencer.matchScore ?? 82;
    score = Math.min(96, Math.max(78, fallback));
  }

  return {
    ...influencer,
    calculatedScore: score,
    bio: influencerBios[influencer.id] ?? "Adfluence creator profile with strong campaign fit.",
    estimatedPriceMin: price.min,
    estimatedPriceMax: price.max,
    creatorSize,
  };
}

export function Recommendations() {
  const [draft, setDraft] = useState<CampaignDraft | null>(null);
  const [invitedIds, setInvitedIds] = useState<string[]>([]);
  const [inviteMessage, setInviteMessage] = useState("");
  const [search, setSearch] = useState("");
  const [nicheFilter, setNicheFilter] = useState("All niches");
  const [platformFilter, setPlatformFilter] = useState("All platforms");
  const [scoreFilter, setScoreFilter] = useState("Any score");
  const [sortBy, setSortBy] = useState("Best match");

  useEffect(() => {
    const rawDraft = localStorage.getItem(draftStorageKey);
    if (rawDraft) {
      try {
        const parsed = JSON.parse(rawDraft) as CampaignDraft;
        setDraft(parsed);
      } catch {
        setDraft(null);
      }
    }

    const rawInvited = localStorage.getItem(invitedStorageKey);
    if (rawInvited) {
      try {
        setInvitedIds(JSON.parse(rawInvited) as string[]);
      } catch {
        setInvitedIds([]);
      }
    }
  }, []);

  const rankedInfluencers = useMemo(
    () =>
      influencers
        .map((influencer) => scoreInfluencer(influencer, draft))
        .sort((a, b) => b.calculatedScore - a.calculatedScore),
    [draft],
  );

  const filteredInfluencers = useMemo(() => {
    const minScore =
      scoreFilter === "80+"
        ? 80
        : scoreFilter === "85+"
          ? 85
          : scoreFilter === "90+"
            ? 90
            : scoreFilter === "95+"
              ? 95
              : 0;

    const searched = rankedInfluencers.filter((influencer) => {
      const bySearch = influencer.name.toLowerCase().includes(search.toLowerCase().trim());
      const byNiche =
        nicheFilter === "All niches" ||
        influencer.niche.toLowerCase() === nicheFilter.toLowerCase();
      const byPlatform =
        platformFilter === "All platforms" || influencer.platforms.includes(platformFilter as never);
      const byScore = influencer.calculatedScore >= minScore;
      return bySearch && byNiche && byPlatform && byScore;
    });

    const sorted = [...searched];
    if (sortBy === "Followers") {
      sorted.sort((a, b) => b.followers - a.followers);
    } else if (sortBy === "Price low to high") {
      sorted.sort((a, b) => a.estimatedPriceMin - b.estimatedPriceMin);
    } else if (sortBy === "Price high to low") {
      sorted.sort((a, b) => b.estimatedPriceMax - a.estimatedPriceMax);
    } else {
      sorted.sort((a, b) => b.calculatedScore - a.calculatedScore);
    }

    return sorted;
  }, [rankedInfluencers, search, nicheFilter, platformFilter, scoreFilter, sortBy]);

  const invitedCreators = useMemo(
    () => rankedInfluencers.filter((influencer) => invitedIds.includes(influencer.id)),
    [rankedInfluencers, invitedIds],
  );

  const onInvite = (id: string, name: string) => {
    if (invitedIds.includes(id)) return;
    const next = [...invitedIds, id];
    setInvitedIds(next);
    localStorage.setItem(invitedStorageKey, JSON.stringify(next));
    setInviteMessage(`Invite sent to ${name}`);
    setTimeout(() => setInviteMessage(""), 2500);
  };

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-white">Recommended Influencers</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            Creators ranked for your campaign based on niche, platform, budget, and audience fit.
          </p>
        </div>
        <Button
          asChild
          className="h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-5 text-white hover:bg-[color:var(--adfluence-orange-hover)]"
        >
          <Link href="/business/campaigns/create">Create New Campaign</Link>
        </Button>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
        <h2 className="text-lg font-semibold text-white">Campaign Context</h2>
        <div className="mt-4 grid gap-3 text-sm text-[#CBD5E1] md:grid-cols-2 xl:grid-cols-4">
          <p>
            Campaign title: <span className="text-white">{draft?.title || "Not set"}</span>
          </p>
          <p>
            Goal: <span className="text-white">{draft?.goal || "Not set"}</span>
          </p>
          <p>
            Niche: <span className="text-white">{draft?.niche || "Not set"}</span>
          </p>
          <p>
            Platforms:{" "}
            <span className="text-white">
              {draft?.platforms?.length ? draft.platforms.join(", ") : "Not set"}
            </span>
          </p>
          <p>
            Budget: <span className="text-white">{draft?.budgetRange || "Not set"}</span>
          </p>
          <p>
            Creator size: <span className="text-white">{draft?.creatorSize || "Not set"}</span>
          </p>
          <p className="md:col-span-2 xl:col-span-2">
            Timeline:{" "}
            <span className="text-white">{formatTimeline(draft?.startDate, draft?.endDate)}</span>
          </p>
        </div>
        <p className="mt-3 text-xs text-[#CBD5E1]">
          Recommendations are generated from your campaign brief.
        </p>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
        <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-5">
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search influencer"
            className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
          />

          <select
            value={nicheFilter}
            onChange={(e) => setNicheFilter(e.target.value)}
            className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
          >
            {nicheOptions.map((option) => (
              <option key={option} value={option} className="bg-[color:var(--adfluence-surface)]">
                {option}
              </option>
            ))}
          </select>

          <select
            value={platformFilter}
            onChange={(e) => setPlatformFilter(e.target.value)}
            className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
          >
            {platformOptions.map((option) => (
              <option key={option} value={option} className="bg-[color:var(--adfluence-surface)]">
                {option}
              </option>
            ))}
          </select>

          <select
            value={scoreFilter}
            onChange={(e) => setScoreFilter(e.target.value)}
            className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
          >
            {scoreOptions.map((option) => (
              <option key={option} value={option} className="bg-[color:var(--adfluence-surface)]">
                {option}
              </option>
            ))}
          </select>

          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
          >
            {sortOptions.map((option) => (
              <option key={option} value={option} className="bg-[color:var(--adfluence-surface)]">
                {option}
              </option>
            ))}
          </select>
        </div>
      </section>

      {inviteMessage ? (
        <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-200">
          {inviteMessage}
        </p>
      ) : null}

      <section className="flex flex-col gap-6 2xl:grid 2xl:grid-cols-[minmax(0,1fr)_320px]">
        <div className="order-2 grid gap-4 sm:grid-cols-2 xl:grid-cols-3 2xl:order-1">
          {filteredInfluencers.map((influencer) => (
            <article
              key={influencer.id}
              className="group rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 shadow-[0_8px_24px_rgba(2,6,23,0.24)] transition-transform duration-200 hover:-translate-y-1"
            >
              <div className="relative h-44 overflow-hidden rounded-xl">
                <Image
                  src={influencer.imageUrl || "/images/influencers/influencer-1.jpg"}
                  alt={influencer.name}
                  fill
                  className="object-cover"
                />
                <span className="absolute right-2 top-2 rounded-full bg-[#0B1A2B]/90 px-2 py-1 text-xs font-medium text-emerald-300">
                  Verified
                </span>
              </div>

              <div className="mt-4 space-y-2">
                <h3 className="text-lg font-semibold text-white">{influencer.name}</h3>
                <p className="text-sm capitalize text-[#CBD5E1]">
                  {influencer.niche} • {influencer.platforms.join(", ")}
                </p>
                <p className="text-sm text-[#CBD5E1]">
                  {influencer.followers.toLocaleString()} followers • {influencer.engagementRate}%
                  engagement
                </p>
                <p className="text-sm text-[#CBD5E1]">
                  Price: ETB {influencer.estimatedPriceMin.toLocaleString()} -{" "}
                  {influencer.estimatedPriceMax.toLocaleString()}
                </p>
                <p className="text-sm text-[#CBD5E1]">{influencer.bio}</p>
              </div>

              <div className="mt-4">
                <div className="mb-1 flex items-center justify-between text-xs">
                  <span className="text-[#CBD5E1]">Match score</span>
                  <span className="font-semibold text-white">{influencer.calculatedScore}%</span>
                </div>
                <div className="h-2 rounded-full bg-white/10">
                  <div
                    className="h-2 rounded-full bg-[color:var(--adfluence-orange)]"
                    style={{ width: `${influencer.calculatedScore}%` }}
                  />
                </div>
              </div>

              <div className="mt-4 grid grid-cols-2 gap-2">
                <Button
                  asChild
                  className="h-10 rounded-xl bg-white/10 text-white hover:bg-white/15"
                >
                  <Link href={`/business/influencers/${influencer.id}`}>View Profile</Link>
                </Button>
                <Button
                  onClick={() => onInvite(influencer.id, influencer.name)}
                  disabled={invitedIds.includes(influencer.id)}
                  className="h-10 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {invitedIds.includes(influencer.id) ? "Invited" : "Invite"}
                </Button>
              </div>
            </article>
          ))}
        </div>

        <aside className="order-1 space-y-4 2xl:order-2">
          <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
            <h3 className="text-lg font-semibold text-white">How matching works</h3>
            <ul className="mt-3 space-y-2 text-sm text-[#CBD5E1]">
              <li>• Niche fit aligns creator content with campaign category.</li>
              <li>• Platform fit prioritizes creators active on selected channels.</li>
              <li>• Budget fit checks creator pricing against your campaign range.</li>
              <li>• Audience fit reflects reach and engagement consistency.</li>
              <li>• Admin verification boosts trusted creator recommendations.</li>
            </ul>
          </div>

          <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
            <h3 className="text-lg font-semibold text-white">Selected / Invited Creators</h3>
            <div className="mt-3 space-y-3">
              {invitedCreators.length === 0 ? (
                <p className="text-sm text-[#CBD5E1]">No creators invited yet.</p>
              ) : (
                invitedCreators.map((creator) => (
                  <div
                    key={creator.id}
                    className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3"
                  >
                    <p className="text-sm font-medium text-white">{creator.name}</p>
                    <p className="mt-1 text-xs text-emerald-300">Status: Invited</p>
                    <Button
                      asChild
                      className="mt-3 h-9 w-full rounded-xl bg-white/10 text-white hover:bg-white/15"
                    >
                      <Link href="/business/messages">Message</Link>
                    </Button>
                  </div>
                ))
              )}
            </div>
          </div>
        </aside>
      </section>
    </div>
  );
}
