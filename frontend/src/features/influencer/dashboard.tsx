"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { ArrowRight, MessageSquare } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  influencerApi,
  type InfluencerDashboardSummaryResponse,
  type ContractListItem,
} from "@/lib/api/influencer";
import { chatApi, type ChatConversationItem } from "@/lib/api/chat";
import { OPTIONAL_PROFILE_ITEMS, REQUIRED_FOR_SUBMIT } from "@/lib/influencer-profile";
import { UserAvatar } from "@/components/common/user-avatar";

function KpiCard({ label, value, hint }: { label: string; value: string; hint?: string }) {
  return (
    <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
      <p className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">{label}</p>
      <p className="mt-2 text-2xl font-bold text-[color:var(--adfluence-orange)]">
        {value}
      </p>
      {hint ? <p className="mt-1 text-xs text-[#CBD5E1]">{hint}</p> : null}
    </article>
  );
}

export function InfluencerDashboard() {
  const [summary, setSummary] = useState<InfluencerDashboardSummaryResponse | null>(null);
  const [conversations, setConversations] = useState<ChatConversationItem[]>([]);
  const [contracts, setContracts] = useState<ContractListItem[]>([]);
  const [loadState, setLoadState] = useState<"loading" | "ready" | "error">("loading");

  const loadDashboard = async () => {
    try {
      const [sumRes, chatRes, contRes] = await Promise.allSettled([
        influencerApi.getDashboardSummary(),
        chatApi.listConversations(),
        influencerApi.listContracts(),
      ]);

      if (sumRes.status !== "fulfilled" || !sumRes.value.data) {
        setLoadState("error");
        return;
      }
      setSummary(sumRes.value.data);
      if (chatRes.status === "fulfilled") {
        setConversations(chatRes.value.data ?? []);
      }
      if (contRes.status === "fulfilled" && contRes.value.data?.items) {
        setContracts(contRes.value.data.items);
      }
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  };

  useEffect(() => {
    void Promise.resolve().then(loadDashboard);
  }, []);

  const kpis = summary?.kpis;
  const completeness = summary?.completeness;
  const earnings = summary?.earnings_summary;
  const pendingInvites = summary?.pending_invites || [];
  const greetingName = summary?.greeting?.creator_name || "Creator";

  const activeContracts = contracts.filter((c) => c.status.toLowerCase() === "active");

  if (loadState === "loading") {
    return <p role="status" className="p-6 text-sm text-[#94A3B8]">Loading your dashboard…</p>;
  }
  if (loadState === "error") {
    return (
      <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
        Your dashboard could not be loaded.{" "}
        <button type="button" className="underline" onClick={() => { setLoadState("loading"); void loadDashboard(); }}>
          Try again
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-white">Welcome back, {greetingName}</h1>
          <p className="mt-1 max-w-xl text-sm text-[#CBD5E1]">
            Manage your creator profile, campaign invites, contracts, and earnings.
          </p>
        </div>
        <Button
          asChild
          className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]"
        >
          <Link href="/influencer/profile">Edit Profile</Link>
        </Button>
      </section>

      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <KpiCard
          label="Campaign Invites"
          value={String(kpis?.pending_invites ?? 0)}
          hint="Awaiting your response"
        />
        <KpiCard
          label="Active Collaborations"
          value={String(kpis?.active_collaborations ?? activeContracts.length)}
        />
        <KpiCard
          label="Earned"
          value={`ETB ${Math.round(kpis?.total_completed_earnings_etb ?? 0).toLocaleString()}`}
          hint="From completed contracts"
        />
        <KpiCard
          label="Profile Score"
          value={`${kpis?.profile_score_percent ?? completeness?.percent ?? 0}%`}
          hint="Complete profiles get better matches"
        />
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:p-6">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h2 className="text-lg font-semibold text-white">Profile strength</h2>
            <p className="mt-1 text-sm text-[#CBD5E1]">A complete profile helps businesses choose you.</p>
          </div>
          <span className="rounded-full border border-white/15 bg-white/5 px-3 py-1 text-sm font-semibold text-white">
            {completeness?.percent ?? 0}%
          </span>
        </div>
        <div className="mt-4 h-2 overflow-hidden rounded-full bg-white/10">
          <div
            className="h-full rounded-full bg-[color:var(--adfluence-action)]"
            style={{ width: `${completeness?.percent ?? 0}%` }}
          />
        </div>
        {completeness?.missing && completeness.missing.length > 0 ? (
          <ul className="mt-4 space-y-2 text-sm text-[#CBD5E1]">
            {completeness.missing.map((item) => (
              <li key={item} className="flex items-center gap-2">
                <span className="text-[color:var(--adfluence-orange)]">○</span>{" "}
                {REQUIRED_FOR_SUBMIT[item] ?? OPTIONAL_PROFILE_ITEMS[item] ?? item.replace(/_/g, " ")}
              </li>
            ))}
          </ul>
        ) : (
          <p className="mt-4 text-sm text-emerald-300">✓ Your profile is complete.</p>
        )}
        <Button
          asChild
          className="mt-4 h-10 rounded-xl bg-white/10 text-white hover:bg-white/15"
        >
          <Link href="/influencer/profile">
            {completeness?.missing?.length ? "Improve profile" : "View profile"}
          </Link>
        </Button>
      </section>

      <section>
        <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
          <h2 className="text-lg font-semibold text-white">Campaign invites</h2>
          <Link
            href="/influencer/invites"
            className="text-sm text-[color:var(--adfluence-orange)] hover:text-[color:var(--adfluence-orange-hover)]"
          >
            View all
          </Link>
        </div>
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {pendingInvites.length === 0 ? (
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-6 text-center text-sm text-[#CBD5E1] col-span-full">
              No pending invites right now.
            </div>
          ) : (
            pendingInvites.slice(0, 3).map((inv) => (
              <article
                key={inv.id}
                className="flex flex-col rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
              >
                <p className="font-semibold text-white">{inv.business_name}</p>
                <p className="text-sm text-[#CBD5E1]">{inv.campaign_title}</p>
                <p className="mt-2 text-xs text-[#94A3B8]">Budget: {inv.budget_summary || "Flexible"}</p>
                {inv.end_date ? (
                  <p className="text-xs text-[#94A3B8]">
                    Deadline: {new Date(inv.end_date).toLocaleDateString()}
                  </p>
                ) : null}
                <span className="mt-2 inline-flex w-fit rounded-full border border-amber-400/30 bg-amber-500/10 px-2 py-0.5 text-xs text-amber-200">
                  Waiting for your reply
                </span>
                {/* Answering happens on the Invites page, which shows the full brief and confirms a decline. */}
                <Button
                  asChild
                  className="mt-4 h-11 w-full rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]"
                >
                  <Link href="/influencer/invites">Review and respond</Link>
                </Button>
              </article>
            ))
          )}
        </div>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:p-6">
        <h2 className="text-lg font-semibold text-white">Earnings overview</h2>
        {/* "Earned" is already a KPI above; this card shows only what is still to be paid. */}
        <div className="mt-4 grid gap-3 sm:grid-cols-2">
          <div className="rounded-xl border border-white/10 bg-[#06111F] p-3">
            <p className="text-xs text-[#CBD5E1]">Pending payout</p>
            <p className="text-lg font-semibold text-amber-200">
              ETB {Math.round(earnings?.pending_payout_etb ?? 0).toLocaleString()}
            </p>
          </div>
          <div className="rounded-xl border border-white/10 bg-[#06111F] p-3">
            <p className="text-xs text-[#CBD5E1]">Next payout</p>
            <p className="text-lg font-semibold text-white">
              {earnings?.next_payout_at
                ? new Date(earnings.next_payout_at).toLocaleDateString()
                : "No payout scheduled"}
            </p>
          </div>
        </div>
      </section>

      <section className="grid gap-4 lg:grid-cols-2">
        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <div className="flex items-center justify-between gap-2">
            <h2 className="text-lg font-semibold text-white">Recent messages</h2>
            <Button asChild variant="ghost" size="sm" className="text-[color:var(--adfluence-orange)] hover:text-white">
              <Link href="/influencer/messages" className="gap-1">
                Open Messages <ArrowRight className="size-4" />
              </Link>
            </Button>
          </div>
          {conversations.length === 0 ? (
            <p className="mt-4 text-xs text-[#CBD5E1]">No active messages yet.</p>
          ) : (
            <ul className="mt-3 space-y-3">
              {conversations.slice(0, 3).map((c) => (
                <li
                  key={c.id}
                  className="flex items-start gap-3 rounded-xl border border-white/10 bg-[#06111F] p-3"
                >
                  <UserAvatar
                    src={c.other_user_avatar}
                    name={c.other_user_name}
                    size="sm"
                  />
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium text-white">{c.other_user_name}</p>
                    <p className="truncate text-xs text-[#CBD5E1]">{c.last_message_preview || "Direct message"}</p>
                  </div>
                  {c.unread_count > 0 ? (
                    <span className="shrink-0 rounded-full bg-[color:var(--adfluence-action)] px-2 py-0.5 text-[10px] font-semibold">
                      {c.unread_count}
                    </span>
                  ) : null}
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <div className="flex items-center justify-between gap-2">
            <h2 className="text-lg font-semibold text-white">Active contracts</h2>
            <MessageSquare className="size-4 text-[color:var(--adfluence-purple-accent)]" aria-hidden />
          </div>
          {activeContracts.length === 0 ? (
            <p className="mt-4 text-xs text-[#CBD5E1]">No active contracts currently.</p>
          ) : (
            <ul className="mt-3 space-y-2">
              {activeContracts.map((c) => (
                <li
                  key={c.id}
                  className="flex items-center justify-between gap-2 rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm"
                >
                  <span className="truncate text-[#CBD5E1]">{c.counterparty_name || c.title}</span>
                  <span className="shrink-0 text-xs text-emerald-300">Active</span>
                </li>
              ))}
            </ul>
          )}
          <Button asChild className="mt-4 h-10 w-full rounded-xl bg-white/10 text-white hover:bg-white/15">
            <Link href="/influencer/contracts">View Contracts</Link>
          </Button>
        </div>
      </section>
    </div>
  );
}
