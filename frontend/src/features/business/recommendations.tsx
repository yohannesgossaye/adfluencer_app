"use client";

import Image from "next/image";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";
import { BadgeCheck, MapPin, Search, SlidersHorizontal, Users, X } from "lucide-react";
import { toast } from "sonner";

import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  businessApi,
  type CampaignListItem,
  type CreatorSearchItem,
  type SearchCreatorsResponse,
} from "@/lib/api/business";
import { usePagedList } from "@/hooks/use-paged-list";
import { businessErrorMessage } from "@/lib/api-errors";
import { cn, getInitials } from "@/lib/utils";
import { formatEtbRange } from "@/lib/contracts";
import { platformLabel } from "@/lib/status";

const PLATFORMS = [
  { value: "", label: "All platforms" },
  { value: "tiktok", label: "TikTok" },
  { value: "instagram", label: "Instagram" },
  { value: "youtube", label: "YouTube" },
  { value: "telegram", label: "Telegram" },
  { value: "facebook", label: "Facebook" },
];

// Same bands as the backend (helper.ClassifyCreatorSize).
const SIZES = [
  { value: "", label: "Any audience size" },
  { value: "nano", label: "Nano · under 10K" },
  { value: "micro", label: "Micro · 10K–70K" },
  { value: "mid", label: "Mid · 70K–110K" },
  { value: "macro", label: "Macro · 110K–500K" },
  { value: "mega", label: "Mega · 500K+" },
];

const SORTS = [
  { value: "best", label: "Best match" },
  { value: "followers", label: "Most followers" },
  { value: "engagement", label: "Highest engagement" },
  { value: "price_low", label: "Price: low to high" },
  { value: "price_high", label: "Price: high to low" },
];

const field =
  "h-11 w-full min-w-0 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";

const compact = new Intl.NumberFormat("en", { notation: "compact", maximumFractionDigits: 1 });

function reach(c: CreatorSearchItem) {
  const parts = [
    c.tiktok_followers ? `TikTok ${compact.format(c.tiktok_followers)}` : "",
    c.instagram_followers ? `Instagram ${compact.format(c.instagram_followers)}` : "",
    c.youtube_subscribers ? `YouTube ${compact.format(c.youtube_subscribers)}` : "",
  ].filter(Boolean);
  return parts.length ? parts.join(" · ") : `${compact.format(c.total_followers)} followers`;
}

const INVITED = new Set(["pending", "accepted"]);

/**
 * Find creators: search and filter every approved creator. Choosing a
 * campaign ranks them by fit (match %) and lets the business invite them.
 */
export function Recommendations() {
  const searchParams = useSearchParams();
  const [campaigns, setCampaigns] = useState<CampaignListItem[]>([]);
  const [categories, setCategories] = useState<Array<{ id: string; name: string; slug: string }>>([]);
  const [campaignId, setCampaignId] = useState(searchParams.get("campaign_id") ?? "");

  const [search, setSearch] = useState("");
  const [q, setQ] = useState("");
  const [niche, setNiche] = useState("");
  const [platform, setPlatform] = useState("");
  const [size, setSize] = useState("");
  const [budget, setBudget] = useState("");
  const [maxPrice, setMaxPrice] = useState<number | undefined>();
  const [sort, setSort] = useState("best");
  const [invitingId, setInvitingId] = useState<string | null>(null);
  // On phones the filters fold away so creators show on the first screen.
  const [filtersOpen, setFiltersOpen] = useState(false);

  useEffect(() => {
    void Promise.allSettled([businessApi.listCampaigns({ page: 1, limit: 50 }), businessApi.listCategories()]).then(([camp, cat]) => {
      if (cat.status === "fulfilled") setCategories(cat.value.data?.items ?? []);
      if (camp.status === "fulfilled") {
        // Only drafts (to preview fit) and live campaigns can be ranked for.
        setCampaigns((camp.value.data?.items ?? []).filter((c) => c.status === "draft" || c.status === "active"));
      }
    });
  }, []);

  // Typing settles for a moment before searching.
  useEffect(() => {
    const t = setTimeout(() => setQ(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);
  useEffect(() => {
    const t = setTimeout(() => {
      const n = Number(budget.replace(/[^\d]/g, ""));
      setMaxPrice(n > 0 ? n : undefined);
    }, 400);
    return () => clearTimeout(t);
  }, [budget]);

  const fetchPage = useCallback(
    async (page: number) =>
      (
        await businessApi.searchCreators({
          q: q || undefined,
          niche: niche || undefined,
          platform: platform || undefined,
          size: size || undefined,
          max_price: maxPrice,
          campaign_id: campaignId || undefined,
          sort,
          page,
        })
      ).data as SearchCreatorsResponse,
    [q, niche, platform, size, maxPrice, campaignId, sort],
  );
  const list = usePagedList(fetchPage);
  const creators: CreatorSearchItem[] = list.items;
  const rankedFor = list.last?.ranked_for;
  const campaign = useMemo(() => campaigns.find((c) => c.id === campaignId), [campaigns, campaignId]);
  const canInvite = campaign?.status === "active";

  const filtered = Boolean(q || niche || platform || size || maxPrice);
  const activeFilters = [campaignId, niche, platform, size, maxPrice].filter(Boolean).length;
  const clearFilters = () => {
    setSearch("");
    setQ("");
    setNiche("");
    setPlatform("");
    setSize("");
    setBudget("");
    setMaxPrice(undefined);
    setSort("best");
  };

  const invite = async (c: CreatorSearchItem) => {
    if (!campaign || !canInvite) return;
    setInvitingId(c.influencer_user_id);
    try {
      await businessApi.createInvite(campaign.id, {
        influencer_user_id: c.influencer_user_id,
        message: `Hello ${c.creator_name}, we'd love to work with you on our campaign "${campaign.title}".`,
      });
      toast.success(`Invitation sent to ${c.creator_name}.`);
      await list.reload();
    } catch (err) {
      toast.error(businessErrorMessage(err, `Could not send the invitation to ${c.creator_name}.`));
    } finally {
      setInvitingId(null);
    }
  };

  return (
    <div className="min-w-0 space-y-5 text-white">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">Find creators</h1>
          <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
            Search every approved creator by name, niche, platform, audience size and price. Choose a campaign to rank them by fit and send invites.
          </p>
        </div>
        <Button asChild className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]">
          <Link href="/business/campaigns/create">Create campaign</Link>
        </Button>
      </section>

      <section aria-label="Search and filters" className="space-y-3 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
        <div className="flex gap-2">
        <div className="relative min-w-0 flex-1">
          <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-[#94A3B8]" aria-hidden />
          <Input
            aria-label="Search creators"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by name, niche, city or topic, e.g. “football”, “Hana”, “Hawassa”"
            className="h-12 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] pl-9 text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
          />
        </div>
          <Button
            type="button"
            variant="outline"
            aria-expanded={filtersOpen}
            aria-controls="creator-filters"
            onClick={() => setFiltersOpen((o) => !o)}
            className="h-12 shrink-0 rounded-xl border-white/15 bg-transparent px-3 text-white hover:bg-white/10 sm:hidden"
          >
            <SlidersHorizontal className="mr-1.5 size-4" aria-hidden />
            Filters{activeFilters ? ` (${activeFilters})` : ""}
          </Button>
        </div>
        <div id="creator-filters" className={cn("gap-3 sm:grid sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6", filtersOpen ? "grid" : "hidden")}>
          <label className="space-y-1 xl:col-span-2">
            <span className="text-xs text-[#94A3B8]">Rank for campaign</span>
            <select value={campaignId} onChange={(e) => setCampaignId(e.target.value)} className={field}>
              <option value="">No campaign: browse all</option>
              {campaigns.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.title}
                  {c.status === "draft" ? " (draft)" : ""}
                </option>
              ))}
            </select>
          </label>
          <label className="space-y-1">
            <span className="text-xs text-[#94A3B8]">Niche</span>
            <select value={niche} onChange={(e) => setNiche(e.target.value)} className={field}>
              <option value="">All niches</option>
              {categories.map((c) => (
                <option key={c.id} value={c.slug}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>
          <label className="space-y-1">
            <span className="text-xs text-[#94A3B8]">Platform</span>
            <select value={platform} onChange={(e) => setPlatform(e.target.value)} className={field}>
              {PLATFORMS.map((p) => (
                <option key={p.value} value={p.value}>
                  {p.label}
                </option>
              ))}
            </select>
          </label>
          <label className="space-y-1">
            <span className="text-xs text-[#94A3B8]">Audience size</span>
            <select value={size} onChange={(e) => setSize(e.target.value)} className={field}>
              {SIZES.map((s) => (
                <option key={s.value} value={s.value}>
                  {s.label}
                </option>
              ))}
            </select>
          </label>
          <label className="space-y-1">
            <span className="text-xs text-[#94A3B8]">Your budget per creator (ETB)</span>
            <Input inputMode="numeric" value={budget} onChange={(e) => setBudget(e.target.value)} placeholder="e.g. 10000" className={cn(field, "h-11")} />
          </label>
        </div>
        <div className="flex flex-wrap items-center justify-between gap-3">
          <p role="status" className="text-sm text-[#CBD5E1]">
            {list.state === "ready" && list.total != null ? (
              <>
                <span className="font-semibold text-white">{list.total}</span> creator{list.total === 1 ? "" : "s"}
                {rankedFor ? (
                  <>
                    {" "}
                    ranked for <span className="text-white">{rankedFor.title}</span>
                  </>
                ) : null}
              </>
            ) : (
              "Searching…"
            )}
          </p>
          <div className="flex items-center gap-2">
            {filtered ? (
              <Button type="button" variant="ghost" onClick={clearFilters} className="h-10 rounded-xl text-white hover:bg-white/10">
                <X className="mr-1 size-4" aria-hidden /> Clear filters
              </Button>
            ) : null}
            <label className="flex items-center gap-2 text-sm text-[#94A3B8]">
              Sort
              <select value={sort} onChange={(e) => setSort(e.target.value)} className={cn(field, "h-10 w-auto")}>
                {SORTS.map((s) => (
                  <option key={s.value} value={s.value}>
                    {s.label}
                  </option>
                ))}
              </select>
            </label>
          </div>
        </div>
        {campaign && campaign.status === "draft" ? (
          <p className="rounded-xl border border-amber-400/30 bg-amber-500/10 px-3 py-2 text-sm text-amber-50">
            “{campaign.title}” is a draft: you can see who fits, but{" "}
            <Link href={`/business/campaigns/${campaign.id}`} className="underline">
              publish it
            </Link>{" "}
            to send invites.
          </p>
        ) : !campaign && campaigns.length > 0 ? (
          <p className="text-xs text-[#94A3B8]">Choose a campaign above to see each creator&apos;s match and invite them.</p>
        ) : null}
      </section>

      <ListStatus
        state={list.state}
        isEmpty={creators.length === 0}
        errorText="Creators could not be loaded."
        onRetry={list.retry}
        emptyText={
          <>
            <Users className="mx-auto size-10 text-[#CBD5E1]/40" aria-hidden />
            <p className="mt-3 text-lg font-medium text-white">No creators match these filters</p>
            <p className="mt-1">Try another niche or platform, or a higher budget.</p>
            {filtered ? (
              <Button type="button" onClick={clearFilters} className="mt-4 h-11 rounded-xl bg-white/10 px-5 text-white hover:bg-white/15">
                Clear filters
              </Button>
            ) : null}
          </>
        }
      >
        <ul className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
          {creators.map((c) => (
            <CreatorCard
              key={c.influencer_user_id}
              creator={c}
              campaignId={campaignId}
              showMatch={Boolean(rankedFor)}
              canInvite={canInvite}
              inviting={invitingId === c.influencer_user_id}
              onInvite={() => void invite(c)}
            />
          ))}
        </ul>
        <LoadMore shown={creators.length} total={list.total} hasMore={list.hasMore} loading={list.loadingMore} onLoadMore={list.loadMore} />
      </ListStatus>
    </div>
  );
}

function CreatorCard({
  creator: c,
  campaignId,
  showMatch,
  canInvite,
  inviting,
  onInvite,
}: {
  creator: CreatorSearchItem;
  campaignId: string;
  showMatch: boolean;
  canInvite: boolean;
  inviting: boolean;
  onInvite: () => void;
}) {
  const invited = c.invite_status ? INVITED.has(c.invite_status) : false;
  const price = c.min_price_etb != null || c.max_price_etb != null ? formatEtbRange(c.min_price_etb, c.max_price_etb) : "Price on request";
  return (
    <li className="flex flex-col rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
      <div className="flex items-start gap-3">
        <div className="relative size-14 shrink-0 overflow-hidden rounded-full border border-white/10 bg-gradient-to-br from-[#4C1D95] to-[color:var(--adfluence-action)]">
          {c.profile_image_url ? (
            <Image src={c.profile_image_url} alt="" fill sizes="56px" className="object-cover" />
          ) : (
            <span className="flex size-full items-center justify-center text-lg font-bold">{getInitials(c.creator_name, "CR")}</span>
          )}
        </div>
        <div className="min-w-0 flex-1">
          <h3 className="flex items-center gap-1.5 font-semibold">
            <span className="truncate">{c.creator_name}</span>
            {c.is_verified ? <BadgeCheck className="size-4 shrink-0 text-sky-300" aria-label="Approved by Adfluence" /> : null}
          </h3>
          <p className="flex items-center gap-1 text-xs text-[#94A3B8]">
            {c.city ? (
              <>
                <MapPin className="size-3" aria-hidden /> {c.city} ·{" "}
              </>
            ) : null}
            <span className="capitalize">{c.size}</span> creator
          </p>
          <div className="mt-1.5 flex flex-wrap gap-1">
            {c.niches.map((n) => (
              <span key={n} className="rounded-full border border-white/15 px-2 py-0.5 text-[11px] text-[#E2E8F0]">
                {n}
              </span>
            ))}
          </div>
        </div>
      </div>

      {c.bio ? <p className="mt-3 line-clamp-2 text-sm text-[#CBD5E1]">{c.bio}</p> : null}

      <dl className="mt-3 space-y-1 text-sm">
        <div className="flex justify-between gap-3">
          <dt className="text-[#94A3B8]">Reach</dt>
          <dd className="truncate text-right">{reach(c)}</dd>
        </div>
        <div className="flex justify-between gap-3">
          <dt className="text-[#94A3B8]">Engagement</dt>
          <dd>{c.engagement_rate ? `${c.engagement_rate}%` : "—"}</dd>
        </div>
        <div className="flex justify-between gap-3">
          <dt className="text-[#94A3B8]">Price</dt>
          <dd className="tabular-nums">{price}</dd>
        </div>
        {c.platforms.length ? (
          <div className="flex justify-between gap-3">
            <dt className="text-[#94A3B8]">Platforms</dt>
            <dd className="truncate text-right">{c.platforms.map(platformLabel).join(", ")}</dd>
          </div>
        ) : null}
      </dl>

      {showMatch && c.match_score != null ? (
        <div className="mt-3">
          <div className="mb-1 flex justify-between text-xs">
            <span className="text-[#CBD5E1]">Campaign match</span>
            <span className="font-semibold">{Math.round(c.match_score)}%</span>
          </div>
          <div className="h-2 rounded-full bg-white/10" aria-hidden>
            <div className="h-2 rounded-full bg-[color:var(--adfluence-action)]" style={{ width: `${Math.min(100, c.match_score)}%` }} />
          </div>
        </div>
      ) : null}

      <div className="mt-auto grid grid-cols-2 gap-2 pt-4">
        <Button asChild className="h-10 rounded-xl bg-white/10 text-white hover:bg-white/15">
          <Link href={`/business/influencers/${c.influencer_user_id}${campaignId ? `?campaign_id=${campaignId}` : ""}`}>View profile</Link>
        </Button>
        <Button
          type="button"
          onClick={onInvite}
          disabled={!canInvite || invited || inviting}
          title={!campaignId ? "Choose a live campaign to invite" : !canInvite ? "Publish the campaign to invite" : undefined}
          className="h-10 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] disabled:opacity-50"
        >
          {inviting ? "Sending…" : invited ? "Invited" : "Invite"}
        </Button>
      </div>
    </li>
  );
}
