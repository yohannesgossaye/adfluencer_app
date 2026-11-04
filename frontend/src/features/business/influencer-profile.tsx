"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { MessageSquare, Send, ShieldCheck } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { businessApi, type InfluencerViewDetail } from "@/lib/api/business";
import { chatApi } from "@/lib/api/chat";
import { apiErrorType, businessErrorMessage } from "@/lib/api-errors";

import { UserAvatar } from "@/components/common/user-avatar";
import { formatEtbRange } from "@/lib/contracts";

type InfluencerProfileProps = {
  influencerId: string;
};

export function InfluencerProfile({ influencerId }: InfluencerProfileProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const urlCampaignId = searchParams.get("campaign_id") || undefined;
  // Invites go to a live campaign only. A campaign_id from a link is preselected
  // when it is live; otherwise the business picks one.
  const [pickedCampaignId, setPickedCampaignId] = useState("");
  const [liveCampaigns, setLiveCampaigns] = useState<Array<{ id: string; title: string }>>([]);
  const urlIsLive = Boolean(urlCampaignId && liveCampaigns.some((c) => c.id === urlCampaignId));
  const campaignId = pickedCampaignId || (urlIsLive ? urlCampaignId : undefined);
  // An invite this creator already has for the chosen campaign (pending or accepted).
  const [existingInvite, setExistingInvite] = useState<{ campaignId: string; status: string } | null>(null);

  const [creator, setCreator] = useState<InfluencerViewDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [inviteSent, setInviteSent] = useState(false);
  const [isInviting, setIsInviting] = useState(false);
  const [loadError, setLoadError] = useState<"not_verified" | "not_found" | "failed" | null>(null);

  useEffect(() => {
    async function loadCreator() {
      try {
        const res = await businessApi.getInfluencerForBusiness(influencerId, urlCampaignId);
        if (res.data) {
          setCreator(res.data);
        }
      } catch (err) {
        const type = apiErrorType(err);
        setLoadError(
          type === "account_not_verified" ? "not_verified" : type === "creator_not_found" ? "not_found" : "failed",
        );
      } finally {
        setIsLoading(false);
      }
    }
    loadCreator();
  }, [influencerId, urlCampaignId]);

  useEffect(() => {
    businessApi
      .listCampaigns({ status: "active", page: 1, limit: 50 })
      .then((res) => setLiveCampaigns((res.data?.items ?? []).map((c) => ({ id: c.id, title: c.title }))))
      .catch(() => setLiveCampaigns([]));
  }, []);

  useEffect(() => {
    if (!campaignId) return;
    let cancelled = false;
    businessApi
      .listInvites(campaignId)
      .then((res) => {
        const inv = res.data?.items?.find(
          (i) => i.influencer_user_id === influencerId && (i.status === "pending" || i.status === "accepted"),
        );
        if (!cancelled) setExistingInvite(inv ? { campaignId, status: inv.status } : null);
      })
      .catch(() => !cancelled && setExistingInvite(null));
    return () => {
      cancelled = true;
    };
  }, [campaignId, influencerId]);
  const alreadyInvited = existingInvite && existingInvite.campaignId === campaignId ? existingInvite.status : null;

  const handleSendInvite = async () => {
    if (!campaignId || !creator) return;
    setIsInviting(true);
    try {
      await businessApi.createInvite(campaignId, {
        influencer_user_id: creator.influencer_user_id,
        message: `Hello ${creator.creator_name}, we'd love to invite you to collaborate with us!`,
      });
      setInviteSent(true);
      toast.success(`Invitation sent to ${creator.creator_name}.`);
    } catch (err) {
      toast.error(businessErrorMessage(err, "The invitation could not be sent."));
    } finally {
      setIsInviting(false);
    }
  };

  const handleStartChat = async () => {
    if (!creator) return;
    try {
      const res = await chatApi.getOrCreateConversation({
        recipient_user_id: creator.influencer_user_id,
        campaign_id: campaignId,
      });
      if (res.data?.id) {
        router.push(`/business/messages?conversation_id=${res.data.id}`);
        return;
      }
      router.push("/business/messages");
    } catch (err) {
      toast.error(businessErrorMessage(err, "The chat could not be opened."));
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center text-[#CBD5E1]">
        Loading creator profile...
      </div>
    );
  }

  if (loadError === "not_verified") {
    return (
      <div className="rounded-2xl border border-amber-400/30 bg-amber-500/10 p-8 text-center text-white">
        <p className="text-lg font-semibold">Creator profiles open once your business is verified</p>
        <p className="mt-1 text-sm text-amber-100/90">
          Adfluence checks every business before it can view and contact creators.
        </p>
        <Button asChild className="mt-4 bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]">
          <Link href="/business/profile">Check my verification status</Link>
        </Button>
      </div>
    );
  }

  if (loadError === "failed") {
    return (
      <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-8 text-center text-sm text-rose-100">
        This creator&apos;s profile could not be loaded. Please try again.
      </div>
    );
  }

  if (!creator) {
    return (
      <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center text-white">
        <p className="text-lg font-semibold">Creator Not Found</p>
        <p className="mt-1 text-sm text-[#CBD5E1]">The creator profile you are looking for does not exist or has been removed.</p>
        <Button asChild className="mt-4 bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]">
          <Link href="/business/recommendations">Back to Find creators</Link>
        </Button>
      </div>
    );
  }

  const primaryNiche = creator.niches?.[0]?.name;
  // Points per criterion as scored by the backend (helper/matching.go).
  const breakdown = creator.score_breakdown;
  const fitRows = breakdown
    ? [
        { label: "Niche fit", value: breakdown.niche, max: 40 },
        { label: "Platform fit", value: breakdown.platform, max: 25 },
        { label: "Budget fit", value: breakdown.budget, max: 20 },
        { label: "Creator size fit", value: breakdown.creator_size, max: 10 },
        { label: "Verified profile", value: breakdown.verified, max: 5 },
      ]
    : [];
  const notProvided = <span className="text-base font-medium text-[#94A3B8]">Not provided</span>;

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-col gap-6 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-6 lg:flex-row lg:items-center lg:justify-between">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
          <UserAvatar
            src={creator.profile_image_url}
            name={creator.creator_name}
            size="2xl"
            fallback="CR"
          />
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <h1 className="text-2xl font-bold text-white">{creator.creator_name}</h1>
              {creator.is_verified ? (
                <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/15 px-2.5 py-0.5 text-xs font-medium text-emerald-300">
                  <ShieldCheck className="size-3.5" />
                  Verified Creator
                </span>
              ) : null}
            </div>
            <p className="mt-1 text-sm text-[#CBD5E1]">
              {[primaryNiche, creator.platforms.join(", "), creator.city].filter(Boolean).join(" • ")}
            </p>
            <p className="mt-1 text-xs text-[#CBD5E1]">
              Full Name: {creator.full_name}
              {creator.languages?.length ? ` · Languages: ${creator.languages.join(", ")}` : ""}
            </p>
          </div>
        </div>

        <div className="flex flex-wrap gap-3">
          <Button
            onClick={handleStartChat}
            variant="outline"
            className="h-11 rounded-xl border-white/15 bg-transparent px-5 text-white hover:bg-white/10"
          >
            <MessageSquare className="mr-2 size-4" />
            Direct Message
          </Button>
          <Button
            onClick={handleSendInvite}
            disabled={!campaignId || inviteSent || isInviting || Boolean(alreadyInvited)}
            className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-6 text-white hover:bg-[color:var(--adfluence-action-hover)] disabled:opacity-60"
          >
            <Send className="mr-2 size-4" />
            {isInviting
              ? "Sending…"
              : inviteSent || alreadyInvited === "pending"
                ? "Invited, waiting for reply"
                : alreadyInvited === "accepted"
                  ? "Accepted your invite"
                  : "Invite to campaign"}
          </Button>
          {(
            liveCampaigns.length > 0 ? (
              <label className="w-full text-xs text-[#CBD5E1]">
                <span className="mb-1 block">Campaign to invite them to</span>
                <select
                  value={campaignId ?? ""}
                  onChange={(e) => {
                    setPickedCampaignId(e.target.value);
                    setInviteSent(false);
                  }}
                  className="h-11 w-full max-w-sm rounded-xl border border-white/10 bg-[#06111F] px-3 text-sm text-white"
                >
                  <option value="">Choose a live campaign…</option>
                  {liveCampaigns.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.title}
                    </option>
                  ))}
                </select>
              </label>
            ) : (
              <p className="w-full text-xs text-[#94A3B8]">
                You have no live campaigns.{" "}
                <Link href="/business/campaigns" className="underline">
                  Publish a campaign
                </Link>{" "}
                to invite creators.
              </p>
            )
          )}
        </div>
      </section>

      <section className="grid gap-3 sm:grid-cols-3">
        <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
          <p className="text-xs text-[#CBD5E1]">Total Followers</p>
          <p className="mt-2 text-2xl font-bold text-white">
            {creator.total_followers.toLocaleString()}
          </p>
        </article>
        <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
          <p className="text-xs text-[#CBD5E1]">Engagement Rate</p>
          <p className="mt-2 text-2xl font-bold text-white">
            {creator.engagement_rate != null ? `${creator.engagement_rate}%` : notProvided}
          </p>
        </article>
        <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
          <p className="text-xs text-[#CBD5E1]">Price Range</p>
          <p className="mt-2 text-xl font-bold text-[color:var(--adfluence-orange)]">
            {creator.min_price_etb != null && creator.max_price_etb != null
              ? formatEtbRange(creator.min_price_etb, creator.max_price_etb)
              : notProvided}
          </p>
        </article>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
        <h2 className="text-lg font-semibold text-white">About Creator</h2>
        <p className="mt-2 text-sm leading-relaxed text-[#CBD5E1]">
          {creator.bio || "This creator hasn't written a bio yet."}
        </p>
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-white">Portfolio &amp; Media Showcase</h2>
        {creator.portfolio_items && creator.portfolio_items.length > 0 ? (
          <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {creator.portfolio_items.map((item) => {
              const mediaUrl = item.file_url || item.media_url;
              if (!mediaUrl) return null;
              return (
                <article
                  key={item.id}
                  className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-3"
                >
                  <div className="relative aspect-video overflow-hidden rounded-xl bg-black/40 border border-white/10">
                    {item.media_type === "video" ? (
                      <video
                        src={mediaUrl}
                        controls
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <Image
                        src={mediaUrl}
                        alt={item.title || "Portfolio item"}
                        fill
                        unoptimized
                        className="object-cover"
                      />
                    )}
                    <span className="absolute left-2 top-2 rounded-full bg-[color:var(--adfluence-surface)]/90 px-2 py-0.5 text-xs capitalize text-[#CBD5E1]">
                      {item.media_type}
                    </span>
                  </div>
                  {item.title ? (
                    <p className="mt-3 text-sm font-medium text-white truncate">{item.title}</p>
                  ) : null}
                  {item.caption ? (
                    <p className="mt-1 text-xs text-[#CBD5E1] line-clamp-2">{item.caption}</p>
                  ) : null}
                </article>
              );
            })}
          </div>
        ) : (
          <div className="rounded-2xl border border-dashed border-white/15 bg-[rgba(255,255,255,0.04)] p-8 text-center">
            <p className="text-sm font-medium text-white">No portfolio or video samples uploaded yet</p>
            <p className="mt-1 text-xs text-[#CBD5E1]">This creator has not added media samples to their showcase.</p>
          </div>
        )}
      </section>

      <section className="grid gap-6 xl:grid-cols-2">
        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <h2 className="text-lg font-semibold text-white">Campaign Fit Analysis</h2>
          {fitRows.length > 0 ? (
            <div className="mt-4 space-y-3">
              {fitRows.map((row) => (
                <div key={row.label}>
                  <div className="mb-1 flex items-center justify-between text-xs">
                    <span className="text-[#CBD5E1]">{row.label}</span>
                    <span className="font-medium text-white">
                      {row.value} / {row.max}
                    </span>
                  </div>
                  <div className="h-2 rounded-full bg-white/10">
                    <div
                      className="h-2 rounded-full bg-[color:var(--adfluence-action)]"
                      style={{ width: `${Math.round((row.value / row.max) * 100)}%` }}
                    />
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="mt-3 text-sm text-[#94A3B8]">
              Open this creator from a campaign&apos;s recommendations to see how well they match it.
            </p>
          )}
        </div>

        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <h2 className="text-lg font-semibold text-white">Platform Channels &amp; Links</h2>
          <div className="mt-4 space-y-3 text-sm">
            {creator.social_links?.tiktok ? (
              <div className="flex items-center justify-between rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
                <span className="text-white font-medium">TikTok</span>
                <span className="text-xs text-[color:var(--adfluence-orange)]">
                  {creator.social_links.tiktok}
                </span>
              </div>
            ) : null}
            {creator.social_links?.instagram ? (
              <div className="flex items-center justify-between rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
                <span className="text-white font-medium">Instagram</span>
                <span className="text-xs text-[color:var(--adfluence-orange)]">
                  {creator.social_links.instagram}
                </span>
              </div>
            ) : null}
            {creator.social_links?.youtube ? (
              <div className="flex items-center justify-between rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
                <span className="text-white font-medium">YouTube</span>
                <span className="text-xs text-[color:var(--adfluence-orange)]">
                  {creator.social_links.youtube}
                </span>
              </div>
            ) : null}
            {creator.social_links?.telegram ? (
              <div className="flex items-center justify-between rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
                <span className="text-white font-medium">Telegram Channel</span>
                <span className="text-xs text-[color:var(--adfluence-orange)]">
                  {creator.social_links.telegram}
                </span>
              </div>
            ) : null}
          </div>
        </div>
      </section>
    </div>
  );
}
