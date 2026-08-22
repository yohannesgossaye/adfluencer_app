"use client";

import Image from "next/image";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { Sparkles, Users, RefreshCw } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  businessApi,
  type CampaignListItem,
  type RecommendationItem,
  type InviteResponse,
} from "@/lib/api/business";
import { getInitials } from "@/lib/utils";

const platformOptions = ["All platforms", "TikTok", "Instagram", "YouTube", "Telegram"] as const;
const scoreOptions = ["Any score", "80+", "85+", "90+", "95+"] as const;
const sortOptions = [
  "Best match",
  "Followers",
  "Price low to high",
  "Price high to low",
] as const;

export function Recommendations() {
  const searchParams = useSearchParams();
  const initialCampaignId = searchParams.get("campaign_id");

  const [campaigns, setCampaigns] = useState<CampaignListItem[]>([]);
  const [categories, setCategories] = useState<Array<{ id: string; name: string; slug: string }>>([]);
  const [selectedCampaignId, setSelectedCampaignId] = useState<string>(initialCampaignId || "");
  const [recommendations, setRecommendations] = useState<RecommendationItem[]>([]);
  const [invites, setInvites] = useState<InviteResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isGenerating, setIsGenerating] = useState(false);
  const [invitingId, setInvitingId] = useState<string | null>(null);
  const [inviteMessage, setInviteMessage] = useState("");

  const [search, setSearch] = useState("");
  const [nicheFilter, setNicheFilter] = useState("All niches");
  const [platformFilter, setPlatformFilter] = useState("All platforms");
  const [scoreFilter, setScoreFilter] = useState("Any score");
  const [sortBy, setSortBy] = useState("Best match");

  // 1. Fetch campaigns and active categories on mount
  useEffect(() => {
    async function loadInitialData() {
      try {
        const [campRes, catRes] = await Promise.allSettled([
          businessApi.listCampaigns({ page: 1, limit: 50 }),
          businessApi.listCategories(),
        ]);

        if (catRes.status === "fulfilled" && catRes.value.data?.items) {
          setCategories(catRes.value.data.items);
        }

        if (campRes.status === "fulfilled" && campRes.value.data?.items) {
          const items = campRes.value.data.items;
          setCampaigns(items);
          if (!selectedCampaignId && items.length > 0) {
            setSelectedCampaignId(items[0].id);
          }
        }
      } catch {
        // Fallback
      } finally {
        setIsLoading(false);
      }
    }
    loadInitialData();
  }, []);

  // 2. Fetch recommendations & invites for selected campaign
  useEffect(() => {
    if (!selectedCampaignId) return;

    async function loadData() {
      setIsLoading(true);
      try {
        const [recRes, invRes] = await Promise.allSettled([
          businessApi.listRecommendations(selectedCampaignId),
          businessApi.listInvites(selectedCampaignId),
        ]);

        if (recRes.status === "fulfilled" && recRes.value.data?.items) {
          let items = recRes.value.data.items;
          // If no recommendations exist yet, try to auto-generate
          if (items.length === 0) {
            await businessApi.generateRecommendations(selectedCampaignId);
            const genRec = await businessApi.listRecommendations(selectedCampaignId);
            items = genRec.data?.items || [];
          }
          setRecommendations(items);
        } else {
          setRecommendations([]);
        }

        if (invRes.status === "fulfilled" && invRes.value.data?.items) {
          setInvites(invRes.value.data.items);
        } else {
          setInvites([]);
        }
      } catch {
        setRecommendations([]);
      } finally {
        setIsLoading(false);
      }
    }

    loadData();
  }, [selectedCampaignId]);

  const handleGenerateRecommendations = async () => {
    if (!selectedCampaignId) return;
    setIsGenerating(true);
    try {
      await businessApi.generateRecommendations(selectedCampaignId);
      const recRes = await businessApi.listRecommendations(selectedCampaignId);
      setRecommendations(recRes.data?.items || []);
      setInviteMessage("AI recommendations refreshed successfully!");
      setTimeout(() => setInviteMessage(""), 3000);
    } catch {
      setInviteMessage("Could not generate recommendations for this campaign.");
      setTimeout(() => setInviteMessage(""), 3000);
    } finally {
      setIsGenerating(false);
    }
  };

  const selectedCampaign = useMemo(
    () => campaigns.find((c) => c.id === selectedCampaignId) ?? campaigns[0],
    [campaigns, selectedCampaignId],
  );

  const invitedInfluencerIds = useMemo(
    () => new Set(invites.map((inv) => inv.influencer_user_id)),
    [invites],
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

    const searched = recommendations.filter((influencer) => {
      const bySearch = (influencer.creator_name || "")
        .toLowerCase()
        .includes(search.toLowerCase().trim());
      const byNiche =
        nicheFilter === "All niches" ||
        (influencer.primary_niche &&
          influencer.primary_niche.toLowerCase() === nicheFilter.toLowerCase());
      const byPlatform =
        platformFilter === "All platforms" ||
        (influencer.platforms &&
          influencer.platforms.some(
            (p) => p.toLowerCase() === platformFilter.toLowerCase(),
          ));
      const byScore = influencer.match_score >= minScore;
      return bySearch && byNiche && byPlatform && byScore;
    });

    const sorted = [...searched];
    if (sortBy === "Followers") {
      sorted.sort((a, b) => (b.total_followers || 0) - (a.total_followers || 0));
    } else if (sortBy === "Price low to high") {
      sorted.sort((a, b) => (a.min_price_etb || 0) - (b.min_price_etb || 0));
    } else if (sortBy === "Price high to low") {
      sorted.sort((a, b) => (b.max_price_etb || 0) - (a.max_price_etb || 0));
    } else {
      sorted.sort((a, b) => b.match_score - a.match_score);
    }

    return sorted;
  }, [recommendations, search, nicheFilter, platformFilter, scoreFilter, sortBy]);

  const onInvite = async (influencerId: string, name: string) => {
    if (!selectedCampaignId || invitedInfluencerIds.has(influencerId)) return;
    setInvitingId(influencerId);
    try {
      const res = await businessApi.createInvite(selectedCampaignId, {
        influencer_user_id: influencerId,
        message: `Hello ${name}, we'd love to collaborate with you on our campaign "${selectedCampaign?.title || "our campaign"}".`,
      });
      const newInvite = res.data;
      if (newInvite) {
        setInvites((prev) => [newInvite, ...prev]);
        setInviteMessage(`Invitation sent to ${name}!`);
        setTimeout(() => setInviteMessage(""), 3000);
      }
    } catch {
      setInviteMessage(`Failed to send invitation to ${name}.`);
      setTimeout(() => setInviteMessage(""), 3000);
    } finally {
      setInvitingId(null);
    }
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
        <div className="flex flex-wrap items-center gap-2">
          {selectedCampaignId ? (
            <Button
              onClick={handleGenerateRecommendations}
              disabled={isGenerating}
              variant="outline"
              className="h-11 rounded-xl border-white/15 bg-transparent px-4 text-white hover:bg-white/10"
            >
              <RefreshCw className={`mr-2 size-4 ${isGenerating ? "animate-spin" : ""}`} />
              Refresh AI Matches
            </Button>
          ) : null}
          <Button
            asChild
            className="h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-5 text-white hover:bg-[color:var(--adfluence-orange-hover)]"
          >
            <Link href="/business/campaigns/create">Create New Campaign</Link>
          </Button>
        </div>
      </section>

      {/* Campaign Selector / Context */}
      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <h2 className="text-lg font-semibold text-white">Select Campaign Context</h2>
          {campaigns.length > 0 ? (
            <select
              value={selectedCampaignId}
              onChange={(e) => setSelectedCampaignId(e.target.value)}
              className="h-10 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm font-medium text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
            >
              {campaigns.map((c) => (
                <option key={c.id} value={c.id} className="bg-[color:var(--adfluence-surface)]">
                  {c.title} ({c.status})
                </option>
              ))}
            </select>
          ) : null}
        </div>

        {selectedCampaign ? (
          <div className="mt-4 grid gap-3 text-sm text-[#CBD5E1] md:grid-cols-2 xl:grid-cols-4">
            <p>
              Campaign title: <span className="text-white font-medium">{selectedCampaign.title}</span>
            </p>
            <p>
              Status: <span className="text-white capitalize">{selectedCampaign.status}</span>
            </p>
            <p>
              Primary Platform:{" "}
              <span className="text-white">{selectedCampaign.primary_platform || "All"}</span>
            </p>
            <p>
              Budget:{" "}
              <span className="text-white">{selectedCampaign.budget_summary || "Standard"}</span>
            </p>
            <p>
              Creators Needed:{" "}
              <span className="text-white">{selectedCampaign.creators_needed}</span>
            </p>
            <p>
              Invites Sent:{" "}
              <span className="text-white">{selectedCampaign.invites_sent}</span>
            </p>
            <p className="md:col-span-2 xl:col-span-2">
              Timeline:{" "}
              <span className="text-white">
                {selectedCampaign.start_date
                  ? `${new Date(selectedCampaign.start_date).toLocaleDateString()} - ${new Date(selectedCampaign.end_date || "").toLocaleDateString()}`
                  : "Flexible"}
              </span>
            </p>
          </div>
        ) : (
          <div className="mt-3 text-sm text-[#CBD5E1]">
            <p>No campaigns created yet. Create a campaign first to get tailored creator recommendations.</p>
          </div>
        )}
      </section>

      {/* Filter Toolbar */}
      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
        <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-5">
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search creator name"
            className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
          />

          <select
            value={nicheFilter}
            onChange={(e) => setNicheFilter(e.target.value)}
            className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
          >
            <option value="All niches" className="bg-[color:var(--adfluence-surface)]">
              All niches
            </option>
            {categories.map((c) => (
              <option key={c.id} value={c.name} className="bg-[color:var(--adfluence-surface)]">
                {c.name}
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
        <div className="order-2 2xl:order-1">
          {filteredInfluencers.length > 0 ? (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
              {filteredInfluencers.map((influencer) => {
                const isInvited = invitedInfluencerIds.has(influencer.influencer_user_id);
                const isInvitingThis = invitingId === influencer.influencer_user_id;

                return (
                  <article
                    key={influencer.influencer_user_id}
                    className="group rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 shadow-[0_8px_24px_rgba(2,6,23,0.24)] transition-transform duration-200 hover:-translate-y-1"
                  >
                    <div className="relative flex h-44 items-center justify-center overflow-hidden rounded-xl bg-gradient-to-br from-[#0B1A2B] via-white/5 to-[color:var(--adfluence-orange)]/20 border border-white/10">
                      {influencer.profile_image_url ? (
                        <Image
                          src={influencer.profile_image_url}
                          alt={influencer.creator_name}
                          fill
                          className="object-cover"
                        />
                      ) : (
                        <div className="flex flex-col items-center justify-center gap-1">
                          <span className="text-3xl font-extrabold text-white uppercase tracking-wider">
                            {getInitials(influencer.creator_name, "CR")}
                          </span>
                          <span className="text-xs text-[#94A3B8]">
                            {influencer.primary_niche || "Creator"}
                          </span>
                        </div>
                      )}
                      {influencer.is_verified ? (
                        <span className="absolute right-2 top-2 rounded-full bg-[#0B1A2B]/90 px-2 py-1 text-xs font-medium text-emerald-300">
                          Verified
                        </span>
                      ) : null}
                    </div>

                    <div className="mt-4 space-y-2">
                      <h3 className="text-lg font-semibold text-white truncate">{influencer.creator_name}</h3>
                      <p className="text-sm capitalize text-[#CBD5E1] truncate">
                        {influencer.primary_niche || "Creator"} •{" "}
                        {influencer.platforms?.length ? influencer.platforms.join(", ") : "Social"}
                      </p>
                      <p className="text-sm text-[#CBD5E1]">
                        {influencer.total_followers.toLocaleString()} followers
                        {influencer.engagement_rate
                          ? ` • ${influencer.engagement_rate}% engagement`
                          : ""}
                      </p>
                      <p className="text-sm text-[#CBD5E1]">
                        Price: ETB {(influencer.min_price_etb || 10000).toLocaleString()} -{" "}
                        {(influencer.max_price_etb || 30000).toLocaleString()}
                      </p>
                    </div>

                    <div className="mt-4">
                      <div className="mb-1 flex items-center justify-between text-xs">
                        <span className="text-[#CBD5E1]">Match score</span>
                        <span className="font-semibold text-white">{influencer.match_score}%</span>
                      </div>
                      <div className="h-2 rounded-full bg-white/10">
                        <div
                          className="h-2 rounded-full bg-[color:var(--adfluence-orange)]"
                          style={{ width: `${influencer.match_score}%` }}
                        />
                      </div>
                    </div>

                    <div className="mt-4 grid grid-cols-2 gap-2">
                      <Button
                        asChild
                        className="h-10 rounded-xl bg-white/10 text-white hover:bg-white/15"
                      >
                        <Link
                          href={`/business/influencers/${influencer.influencer_user_id}${
                            selectedCampaignId ? `?campaign_id=${selectedCampaignId}` : ""
                          }`}
                        >
                          View Profile
                        </Link>
                      </Button>
                      <Button
                        onClick={() =>
                          onInvite(influencer.influencer_user_id, influencer.creator_name)
                        }
                        disabled={isInvited || isInvitingThis || !selectedCampaignId}
                        className="h-10 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] disabled:cursor-not-allowed disabled:opacity-60"
                      >
                        {isInvitingThis
                          ? "Sending..."
                          : isInvited
                            ? "Invited"
                            : "Invite"}
                      </Button>
                    </div>
                  </article>
                );
              })}
            </div>
          ) : (
            <div className="rounded-2xl border border-white/10 bg-white/5 p-12 text-center">
              <Users className="mx-auto size-12 text-[#CBD5E1]/40" />
              <p className="mt-3 text-lg font-medium text-white">No creator matches found</p>
              <p className="mt-1 text-sm text-[#CBD5E1]">
                {campaigns.length === 0
                  ? "Create a campaign first to generate AI creator matches."
                  : "Try adjusting your filters or click Refresh AI Matches to generate recommendations."}
              </p>
              {selectedCampaignId ? (
                <Button
                  onClick={handleGenerateRecommendations}
                  disabled={isGenerating}
                  className="mt-4 h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-5 text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                >
                  <Sparkles className="mr-2 size-4" />
                  Generate Matches for This Campaign
                </Button>
              ) : null}
            </div>
          )}
        </div>

        {/* Right Sidebar: Matching details & Live Sent Invites */}
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
              {invites.length === 0 ? (
                <p className="text-sm text-[#CBD5E1]">No creators invited yet.</p>
              ) : (
                invites.map((invite) => (
                  <div
                    key={invite.id}
                    className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3"
                  >
                    <p className="text-sm font-medium text-white">
                      {invite.creator_name || "Creator"}
                    </p>
                    <p className="mt-1 text-xs text-emerald-300 capitalize">
                      Status: {invite.status}
                    </p>
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
