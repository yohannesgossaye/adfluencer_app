"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import axios from "axios";
import { ArrowLeft, Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { businessApi, type CampaignDetail, type InviteResponse } from "@/lib/api/business";
import { deliverableLabel, formatDate, formatEtb } from "@/lib/contracts";
import { cn } from "@/lib/utils";

const STATUS: Record<string, { label: string; className: string; hint: string }> = {
  draft: { label: "Draft", className: "border-amber-400/30 bg-amber-500/10 text-amber-200", hint: "Only you can see this campaign. Publish it to invite creators." },
  active: { label: "Live", className: "border-emerald-400/30 bg-emerald-500/10 text-emerald-200", hint: "Creators you invite can see and accept this campaign." },
  paused: { label: "Paused by Adfluence", className: "border-rose-400/30 bg-rose-500/10 text-rose-200", hint: "Our team paused this campaign. Contact support if you think this is a mistake." },
  completed: { label: "Completed", className: "border-sky-400/30 bg-sky-500/10 text-sky-200", hint: "This campaign is finished." },
  cancelled: { label: "Cancelled", className: "border-white/15 bg-white/5 text-[#CBD5E1]", hint: "This campaign was cancelled." },
};

const INVITE_STATUS: Record<string, string> = {
  pending: "Waiting for reply",
  accepted: "Accepted",
  rejected: "Declined",
  cancelled: "Cancelled",
};

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { message?: string; error?: { type?: string } } | undefined;
    switch (data?.error?.type) {
      case "account_not_verified":
        return "Your business must be verified before publishing. Check your business profile.";
      case "profile_incomplete":
        return "Add a goal and a budget to the campaign before publishing.";
      case "invalid_campaign_status":
        return "This campaign’s status doesn’t allow that action.";
    }
    return data?.message || (!err.response ? "You appear to be offline." : fallback);
  }
  return fallback;
}

function Info({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-white/10 bg-[rgba(255,255,255,0.04)] px-3 py-2">
      <p className="text-xs text-[#94A3B8]">{label}</p>
      <p className="mt-0.5 text-sm text-white">{children || "—"}</p>
    </div>
  );
}

export function CampaignDetailPage({ campaignId }: { campaignId: string }) {
  const justCreated = useSearchParams().get("created") === "1";
  const [campaign, setCampaign] = useState<CampaignDetail | null>(null);
  const [invites, setInvites] = useState<InviteResponse[]>([]);
  const [canPublish, setCanPublish] = useState<boolean | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [confirm, setConfirm] = useState<"publish" | "cancel" | null>(null);
  const [confirmInvite, setConfirmInvite] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    try {
      const [c, inv, profile] = await Promise.all([
        businessApi.getCampaign(campaignId),
        businessApi.listInvites(campaignId).catch(() => null),
        businessApi.getProfile().catch(() => null),
      ]);
      setCampaign(c.data ?? null);
      setInvites(inv?.data?.items ?? []);
      setCanPublish(profile?.data?.marketplace_access.can_publish_campaigns ?? null);
      setLoadError(null);
    } catch (err) {
      setLoadError(apiError(err, "Could not load this campaign."));
    }
  }, [campaignId]);

  useEffect(() => {
    void Promise.resolve().then(load);
  }, [load]);

  const run = async () => {
    if (!campaign || !confirm) return;
    setBusy(true);
    try {
      if (confirm === "publish") {
        await businessApi.publishCampaign(campaign.id);
        toast.success("Campaign published. You can now invite creators.");
      } else {
        await businessApi.cancelCampaign(campaign.id);
        toast.success("Campaign cancelled. Pending invites were withdrawn.");
      }
      setConfirm(null);
      await load();
    } catch (err) {
      toast.error(apiError(err, "The action could not be completed."));
    } finally {
      setBusy(false);
    }
  };

  const cancelInvite = async (invite: InviteResponse) => {
    setBusy(true);
    try {
      await businessApi.cancelInvite(campaignId, invite.id);
      toast.success(`Invite to ${invite.creator_name || "the creator"} withdrawn.`);
      setConfirmInvite(null);
      await load();
    } catch (err) {
      toast.error(apiError(err, "Could not withdraw the invite."));
    } finally {
      setBusy(false);
    }
  };

  if (loadError) {
    return (
      <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
        {loadError}{" "}
        <Link href="/business/campaigns" className="underline">Back to campaigns</Link>
      </div>
    );
  }
  if (!campaign) {
    return <p className="flex items-center gap-2 py-10 text-sm text-[#94A3B8]"><Loader2 className="size-4 animate-spin" aria-hidden /> Loading campaign…</p>;
  }

  const status = STATUS[campaign.status] ?? { label: campaign.status, className: "border-white/15 text-white", hint: "" };
  const editable = ["draft", "active", "paused"].includes(campaign.status);
  const budget =
    campaign.budget_min_etb != null && campaign.budget_max_etb != null
      ? `${formatEtb(campaign.budget_min_etb)} – ${formatEtb(campaign.budget_max_etb)}`
      : "";

  return (
    <div className="space-y-6 text-white">
      <Link href="/business/campaigns" className="inline-flex items-center gap-1 text-sm text-[#CBD5E1] hover:text-white">
        <ArrowLeft className="size-4" aria-hidden /> All campaigns
      </Link>

      {justCreated && campaign.status === "draft" ? (
        <p role="status" className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">
          Draft saved. Review it below, then publish when you’re ready to invite creators.
        </p>
      ) : null}

      <section className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div className="min-w-0">
            <div className="flex flex-wrap items-center gap-2">
              <h1 className="text-2xl font-semibold">{campaign.title}</h1>
              <span className={cn("rounded-full border px-2 py-0.5 text-xs font-medium", status.className)}>{status.label}</span>
            </div>
            <p className="mt-1 text-sm text-[#CBD5E1]">{status.hint}</p>
          </div>
        </div>

        {confirm ? (
          <div className="space-y-2 rounded-xl border border-white/15 bg-black/20 p-3">
            <p className="text-sm">
              {confirm === "publish"
                ? "Publish this campaign? It goes live and you can invite creators to it."
                : "Cancel this campaign? Pending invites are withdrawn and it can’t be reopened."}
            </p>
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                disabled={busy}
                onClick={() => void run()}
                className={cn("h-10 rounded-xl text-white", confirm === "publish" ? "bg-[color:var(--adfluence-orange)] hover:bg-[color:var(--adfluence-orange-hover)]" : "bg-rose-600 hover:bg-rose-700")}
              >
                {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
                {confirm === "publish" ? "Publish campaign" : "Cancel campaign"}
              </Button>
              <Button type="button" variant="ghost" disabled={busy} onClick={() => setConfirm(null)} className="h-10 rounded-xl text-white hover:bg-white/10">
                Not now
              </Button>
            </div>
          </div>
        ) : (
          <div className="flex flex-wrap items-center gap-2">
            {campaign.status === "draft" ? (
              canPublish === false ? (
                <p className="rounded-xl border border-amber-400/30 bg-amber-500/10 px-3 py-2 text-sm text-amber-100">
                  Publishing unlocks once your business is verified.{" "}
                  <Link href="/business/profile" className="underline">Check verification status</Link>
                </p>
              ) : (
                <Button type="button" onClick={() => setConfirm("publish")} className="h-10 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]">
                  Publish
                </Button>
              )
            ) : null}
            {campaign.status === "active" ? (
              <Button asChild className="h-10 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]">
                <Link href={`/business/recommendations?campaign_id=${campaign.id}`}>Find creators</Link>
              </Button>
            ) : null}
            {editable ? (
              <Button asChild variant="outline" className="h-10 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10">
                <Link href={`/business/campaigns/${campaign.id}/edit`}>Edit</Link>
              </Button>
            ) : null}
            {editable ? (
              <Button type="button" variant="outline" onClick={() => setConfirm("cancel")} className="h-10 rounded-xl border-rose-400/40 bg-transparent text-rose-200 hover:bg-rose-500/10">
                Cancel campaign
              </Button>
            ) : null}
          </div>
        )}

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <Info label="Goal">{campaign.goal}</Info>
          <Info label="Niche">{campaign.category?.name}</Info>
          <Info label="Location">{campaign.city}</Info>
          <Info label="Platforms">{campaign.platforms.map((p) => p.charAt(0).toUpperCase() + p.slice(1)).join(", ")}</Info>
          <Info label="Budget">{budget}</Info>
          <Info label="Creators needed">{campaign.number_of_creators_needed}</Info>
          <Info label="Creator size">{campaign.creator_size_preference}</Info>
          <Info label="Dates">{`${formatDate(campaign.start_date)} → ${formatDate(campaign.end_date)}`}</Info>
          <Info label="Deliverables">{campaign.deliverables.map((d) => `${d.quantity} × ${deliverableLabel(d.type)}`).join(", ")}</Info>
        </div>
        <div>
          <p className="text-xs text-[#94A3B8]">Description</p>
          <p className="mt-1 text-sm whitespace-pre-wrap text-[#E2E8F0]">{campaign.description}</p>
        </div>
      </section>

      <section className="space-y-3 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <h2 className="text-lg font-semibold">Invited creators</h2>
          <p className="text-xs text-[#94A3B8]">
            {campaign.invites_accepted} accepted · {campaign.invites_sent} sent · {campaign.number_of_creators_needed} needed
          </p>
        </div>
        {invites.length === 0 ? (
          <p className="text-sm text-[#CBD5E1]">
            {campaign.status === "active"
              ? "No invites yet. Use Find creators to invite matched creators."
              : campaign.status === "draft"
                ? "Invites can be sent once the campaign is published."
                : "No creators were invited to this campaign."}
          </p>
        ) : (
          <ul className="space-y-2">
            {invites.map((inv) => (
              <li key={inv.id} className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-white/10 bg-[rgba(255,255,255,0.04)] px-3 py-2">
                <div className="min-w-0">
                  <Link href={`/business/influencers/${inv.influencer_user_id}?campaign_id=${campaign.id}`} className="font-medium text-white hover:underline">
                    {inv.creator_name || "Creator"}
                  </Link>
                  <p className="text-xs text-[#94A3B8]">{INVITE_STATUS[inv.status] ?? inv.status} · sent {formatDate(inv.created_at)}</p>
                </div>
                {inv.status === "pending" ? (
                  confirmInvite === inv.id ? (
                    <span className="flex gap-2">
                      <Button type="button" size="sm" disabled={busy} onClick={() => void cancelInvite(inv)} className="h-9 rounded-lg bg-rose-600 text-white hover:bg-rose-700">
                        Withdraw invite
                      </Button>
                      <Button type="button" size="sm" variant="ghost" disabled={busy} onClick={() => setConfirmInvite(null)} className="h-9 rounded-lg text-white hover:bg-white/10">
                        Keep
                      </Button>
                    </span>
                  ) : (
                    <Button type="button" size="sm" variant="outline" onClick={() => setConfirmInvite(inv.id)} className="h-9 rounded-lg border-white/15 bg-transparent text-[#CBD5E1] hover:bg-white/10">
                      Withdraw
                    </Button>
                  )
                ) : inv.status === "accepted" && campaign.status === "active" ? (
                  <Button asChild size="sm" className="h-9 rounded-lg bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]">
                    <Link href={`/business/contracts?campaign_id=${campaign.id}&invite_id=${inv.id}`}>Send contract</Link>
                  </Button>
                ) : null}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
