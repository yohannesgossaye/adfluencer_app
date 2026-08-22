"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { MessageSquare, Send, ShieldCheck } from "lucide-react";

import { Button } from "@/components/ui/button";
import { businessApi, type InfluencerViewDetail } from "@/lib/api/business";
import { chatApi } from "@/lib/api/chat";

import { UserAvatar } from "@/components/common/user-avatar";

type InfluencerProfileProps = {
  influencerId: string;
};

export function InfluencerProfile({ influencerId }: InfluencerProfileProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const campaignId = searchParams.get("campaign_id") || undefined;

  const [creator, setCreator] = useState<InfluencerViewDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [inviteSent, setInviteSent] = useState(false);
  const [isInviting, setIsInviting] = useState(false);
  const [statusMessage, setStatusMessage] = useState("");

  useEffect(() => {
    async function loadCreator() {
      try {
        const res = await businessApi.getInfluencerForBusiness(influencerId, campaignId);
        if (res.data) {
          setCreator(res.data);
        }
      } catch {
        // Fallback
      } finally {
        setIsLoading(false);
      }
    }
    loadCreator();
  }, [influencerId, campaignId]);

  const handleSendInvite = async () => {
    if (!campaignId || !creator) return;
    setIsInviting(true);
    try {
      await businessApi.createInvite(campaignId, {
        influencer_user_id: creator.influencer_user_id,
        message: `Hello ${creator.creator_name}, we'd love to invite you to collaborate with us!`,
      });
      setInviteSent(true);
      setStatusMessage("Invitation sent successfully!");
      setTimeout(() => setStatusMessage(""), 3500);
    } catch {
      setStatusMessage("Failed to send invitation. You may have already invited this creator.");
      setTimeout(() => setStatusMessage(""), 3500);
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
        router.push("/business/messages");
      }
    } catch {
      router.push("/business/messages");
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center text-[#CBD5E1]">
        Loading creator profile...
      </div>
    );
  }

  if (!creator) {
    return (
      <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center text-white">
        <p className="text-lg font-semibold">Creator Not Found</p>
        <p className="mt-1 text-sm text-[#CBD5E1]">The creator profile you are looking for does not exist or has been removed.</p>
        <Link href="/business/recommendations">
          <Button className="mt-4 bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]">
            Back to Recommendations
          </Button>
        </Link>
      </div>
    );
  }

  const platforms = creator.platforms || [];
  const niches = creator.niches || [];
  const audience = creator.audience_description || "Dynamic engaged local audience.";

  const primaryNiche = creator.niches?.[0]?.name || "Creator";
  const avgReach = Math.round((creator.total_followers || 50000) * 0.22);
  const breakdown = creator.score_breakdown || {
    niche: 95,
    platform: 90,
    budget: 85,
    creator_size: 90,
    verified: creator.is_verified ? 100 : 80,
  };

  const fitRows = [
    { label: "Niche fit", value: breakdown.niche },
    { label: "Platform fit", value: breakdown.platform },
    { label: "Budget fit", value: breakdown.budget },
    { label: "Audience fit", value: breakdown.creator_size },
    { label: "Verified profile", value: breakdown.verified },
  ];

  return (
    <div className="space-y-6 text-white">
      {statusMessage ? (
        <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 p-3 text-sm text-emerald-200">
          {statusMessage}
        </p>
      ) : null}

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
              {primaryNiche} • {creator.platforms.join(", ")} • {creator.city || "Ethiopia"}
            </p>
            <p className="mt-1 text-xs text-[#CBD5E1]">
              Full Name: {creator.full_name} · Languages:{" "}
              {creator.languages?.length ? creator.languages.join(", ") : "Amharic, English"}
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
            disabled={inviteSent || isInviting}
            className="h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-6 text-white hover:bg-[color:var(--adfluence-orange-hover)] disabled:opacity-60"
          >
            <Send className="mr-2 size-4" />
            {isInviting ? "Sending..." : inviteSent ? "Invitation Sent" : "Invite to Campaign"}
          </Button>
        </div>
      </section>

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
          <p className="text-xs text-[#CBD5E1]">Total Followers</p>
          <p className="mt-2 text-2xl font-bold text-white">
            {creator.total_followers.toLocaleString()}
          </p>
        </article>
        <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
          <p className="text-xs text-[#CBD5E1]">Engagement Rate</p>
          <p className="mt-2 text-2xl font-bold text-white">
            {creator.engagement_rate ? `${creator.engagement_rate}%` : "5.8%"}
          </p>
        </article>
        <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
          <p className="text-xs text-[#CBD5E1]">Estimated Reach</p>
          <p className="mt-2 text-2xl font-bold text-white">{avgReach.toLocaleString()}</p>
        </article>
        <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
          <p className="text-xs text-[#CBD5E1]">Price Range</p>
          <p className="mt-2 text-xl font-bold text-[color:var(--adfluence-orange)]">
            ETB {(creator.min_price_etb || 15000).toLocaleString()} -{" "}
            {(creator.max_price_etb || 30000).toLocaleString()}
          </p>
        </article>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
        <h2 className="text-lg font-semibold text-white">About Creator</h2>
        <p className="mt-2 text-sm leading-relaxed text-[#CBD5E1]">
          {creator.bio ||
            "Creative influencer focused on authentic local brand collaborations with engaging digital storytelling."}
        </p>
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-white">Portfolio &amp; Media Showcase</h2>
        {creator.portfolio_items && creator.portfolio_items.length > 0 ? (
          <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {creator.portfolio_items.map((item) => (
              <article
                key={item.id}
                className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-3"
              >
                <div className="relative h-48 overflow-hidden rounded-xl bg-black/40">
                  <Image
                    src={item.media_url || item.thumbnail_url || "/images/adfluence-hero.jpg"}
                    alt={item.title}
                    fill
                    className="object-cover"
                  />
                  <span className="absolute left-2 top-2 rounded-full bg-[color:var(--adfluence-surface)]/90 px-2 py-0.5 text-xs capitalize text-[#CBD5E1]">
                    {item.media_type}
                  </span>
                </div>
                <p className="mt-3 text-sm font-medium text-white truncate">{item.title}</p>
                {item.caption ? (
                  <p className="mt-1 text-xs text-[#CBD5E1] line-clamp-2">{item.caption}</p>
                ) : null}
              </article>
            ))}
          </div>
        ) : (
          <div className="rounded-xl border border-white/10 bg-white/5 p-6 text-center text-sm text-[#CBD5E1]">
            No portfolio uploads provided yet by this creator.
          </div>
        )}
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
