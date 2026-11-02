"use client";

import Image from "next/image";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";
import { Sparkles, Users, RefreshCw } from "lucide-react";

import { toast } from "sonner";

import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  businessApi,
  type CampaignListItem,
  type RecommendationItem,
  type InviteResponse,
} from "@/lib/api/business";
import { usePagedList, type LoadState } from "@/hooks/use-paged-list";
import { apiErrorType, businessErrorMessage } from "@/lib/api-errors";
import { getInitials } from "@/lib/utils";

const platformOptions = ["All platforms", "TikTok", "Instagram", "YouTube", "Telegram"] as const;
const scoreOptions = ["Any score", "80+", "85+", "90+", "95+"] as const;
const sortOptions = [
  "Best match",
  "Followers",
  "Price low to high",
  "Price high to low",
] as const;

// Filters run on the server so they cover every match, not just the first page.
const MIN_SCORE: Record<string, number | undefined> = { "80+": 80, "85+": 85, "90+": 90, "95+": 95 };
const SORT: Record<string, string> = {
  "Best match": "match",
  Followers: "followers",
  "Price low to high": "price_low",
  "Price high to low": "price_high",
};

export function Recommendations() {
  const searchParams = useSearchParams();
  const initialCampaignId = searchParams.get("campaign_id");

  const [campaigns, setCampaigns] = useState<CampaignListItem[]>([]);
  const [categories, setCategories] = useState<Array<{ id: string; name: string; slug: string }>>([]);
  const [selectedCampaignId, setSelectedCampaignId] = useState<string>(initialCampaignId || "");
  const [invites, setInvites] = useState<InviteResponse[]>([]);
  const [campaignsState, setCampaignsState] = useState<LoadState>("loading");
  const [isGenerating, setIsGenerating] = useState(false);
  const [invitingId, setInvitingId] = useState<string | null>(null);

  const [search, setSearch] = useState("");
  const [query, setQuery] = useState("");
  const [nicheFilter, setNicheFilter] = useState("All niches");
  const [platformFilter, setPlatformFilter] = useState("All platforms");
  const [scoreFilter, setScoreFilter] = useState("Any score");
  const [sortBy, setSortBy] = useState("Best match");

  // 1. Campaigns (for the picker) and categories (for the niche filter).
  const loadInitialData = useCallback(async () => {
    const [campRes, catRes] = await Promise.allSettled([
      businessApi.listCampaigns({ page: 1, limit: 50 }),
      businessApi.listCategories(),
    ]);
    if (catRes.status === "fulfilled" && catRes.value.data?.items) {
      setCategories(catRes.value.data.items);
    }
    if (campRes.status === "rejected") {
      setCampaignsState("error");
      return;
    }
    // Only drafts (preview) and live campaigns can have recommendations.
    const items = (campRes.value.data?.items ?? []).filter((c) => c.status === "draft" || c.status === "active");
    setCampaigns(items);
    // A campaign_id from a link is used only if it can have matches (draft or live).
    setSelectedCampaignId((current) => (items.some((c) => c.id === current) ? current : items[0]?.id || ""));
    setCampaignsState("ready");
  }, []);

  useEffect(() => {
    void Promise.resolve().then(loadInitialData);
  }, [loadInitialData]);

  useEffect(() => {
    const t = setTimeout(() => setQuery(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  // 2. Matches for the selected campaign, filtered and sorted by the server.
  const fetchPage = useCallback(
    async (page: number) => {
      if (!selectedCampaignId || campaignsState !== "ready" || !campaigns.some((c) => c.id === selectedCampaignId)) {
        return { items: [] as RecommendationItem[], pagination: { total_docs: 0 } };
      }
      const params = {
        search: query || undefined,
        niche: nicheFilter === "All niches" ? undefined : nicheFilter,
        platform: platformFilter === "All platforms" ? undefined : platformFilter.toLowerCase(),
        min_score: MIN_SCORE[scoreFilter],
        sort: SORT[sortBy],
        page,
      };
      let res = (await businessApi.listRecommendations(selectedCampaignId, params)).data;
      const filtered = Boolean(params.search || params.niche || params.platform || params.min_score);
      // A new campaign has no matches yet: generate them once, then list again.
      if (page === 1 && !filtered && (res?.items?.length ?? 0) === 0) {
        await businessApi.generateRecommendations(selectedCampaignId);
        res = (await businessApi.listRecommendations(selectedCampaignId, params)).data;
      }
      return res;
    },
    [selectedCampaignId, campaignsState, campaigns, query, nicheFilter, platformFilter, scoreFilter, sortBy],
  );
  const list = usePagedList(fetchPage);
  const recommendations: RecommendationItem[] = list.items;
  const notVerified = list.state === "error" && apiErrorType(list.error) === "account_not_verified";

  // 3. Invites already sent for this campaign (to mark "Invited").
  useEffect(() => {
    if (!selectedCampaignId) return;
    let cancelled = false;
    businessApi
      .listInvites(selectedCampaignId)
      .then((res) => !cancelled && setInvites(res.data?.items ?? []))
      .catch(() => !cancelled && setInvites([]));
    return () => {
      cancelled = true;
    };
  }, [selectedCampaignId]);

  const handleGenerateRecommendations = async () => {
    if (!selectedCampaignId) return;
    setIsGenerating(true);
    try {
      await businessApi.generateRecommendations(selectedCampaignId);
      await list.reload();
      toast.success("Matches refreshed.");
    } catch (err) {
      toast.error(businessErrorMessage(err, "Could not refresh matches for this campaign."));
    } finally {
      setIsGenerating(false);
    }
  };

  const selectedCampaign = useMemo(
    () => campaigns.find((c) => c.id === selectedCampaignId) ?? campaigns[0],
    [campaigns, selectedCampaignId],
  );

  // Declined or withdrawn invites can be sent again (the server reopens them).
  const invitedInfluencerIds = useMemo(
    () =>
      new Set(
        invites.filter((inv) => inv.status === "pending" || inv.status === "accepted").map((inv) => inv.influencer_user_id),
      ),
    [invites],
  );
  const isDraft = selectedCampaign?.status === "draft";

  const onInvite = async (influencerId: string, name: string) => {
    if (!selectedCampaignId || invitedInfluencerIds.has(influencerId)) return;
    setInvitingId(influencerId);
    try {
      const res = await businessApi.createInvite(selectedCampaignId, {
        influencer_user_id: influencerId,
        message: `Hello ${name}, we'd love to collaborate with you on our campaign "${selectedCampaign?.title || "our campaign"}".`,
      });
      if (res.data?.id) {
        const inv = await businessApi.listInvites(selectedCampaignId);
        setInvites(inv.data?.items ?? []);
        toast.success(`Invitation sent to ${name}.`);
      }
    } catch (err) {
      toast.error(businessErrorMessage(err, `Could not send the invitation to ${name}.`));
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
            className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]"
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

        {isDraft && selectedCampaign ? (
          <p role="status" className="mt-4 rounded-xl border border-amber-400/30 bg-amber-500/10 p-3 text-sm text-amber-50">
            This campaign is still a draft. You can preview matches, but{" "}
            <Link href={`/business/campaigns/${selectedCampaign.id}`} className="underline">
              publish it
            </Link>{" "}
            before inviting creators.
          </p>
        ) : null}
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
            aria-label="Search creator name"
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

      {campaignsState === "error" ? (
        <p role="alert" className="rounded-xl border border-rose-400/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-100">
          Your campaigns could not be loaded.{" "}
          <button
            type="button"
            className="underline"
            onClick={() => {
              setCampaignsState("loading");
              void loadInitialData();
            }}
          >
            Try again
          </button>
        </p>
      ) : null}

      <section className="flex flex-col gap-6 2xl:grid 2xl:grid-cols-[minmax(0,1fr)_320px]">
        <div className="order-2 2xl:order-1">
          {notVerified ? (
            <p role="alert" className="rounded-xl border border-amber-400/30 bg-amber-500/10 px-3 py-2 text-sm text-amber-100">
              Recommended creators appear once your business is verified. Check your{" "}
              <Link href="/business/profile" className="underline">
                business profile
              </Link>{" "}
              for the status.
            </p>
          ) : list.state !== "ready" || campaignsState === "loading" ? (
            <ListStatus
              state={campaignsState === "loading" ? "loading" : list.state}
              isEmpty={false}
              emptyText=""
              errorText="Matches could not be loaded."
              onRetry={list.retry}
            >
              {null}
            </ListStatus>
          ) : recommendations.length > 0 ? (
            <>
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
              {recommendations.map((influencer) => {
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
                        {influencer.min_price_etb != null && influencer.max_price_etb != null
                          ? `Price: ETB ${influencer.min_price_etb.toLocaleString()} – ${influencer.max_price_etb.toLocaleString()}`
                          : "Price not set"}
                      </p>
                    </div>

                    <div className="mt-4">
                      <div className="mb-1 flex items-center justify-between text-xs">
                        <span className="text-[#CBD5E1]">Match score</span>
                        <span className="font-semibold text-white">{influencer.match_score}%</span>
                      </div>
                      <div className="h-2 rounded-full bg-white/10">
                        <div
                          className="h-2 rounded-full bg-[color:var(--adfluence-action)]"
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
                        disabled={isInvited || isInvitingThis || !selectedCampaignId || isDraft}
                        title={isDraft ? "Publish this campaign to invite creators" : undefined}
                        className="h-10 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] disabled:cursor-not-allowed disabled:opacity-60"
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
            <LoadMore
              shown={recommendations.length}
              total={list.total}
              hasMore={list.hasMore}
              loading={list.loadingMore}
              onLoadMore={list.loadMore}
            />
            </>
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
                  className="mt-4 h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]"
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
