"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { ArrowRight, MessageSquare } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  initialInfluencerInvites,
  loadInfluencerInvites,
  saveInfluencerInvites,
  type InfluencerInvite,
} from "@/data/influencer-invites-mock";
import { initialCreatorConversations } from "@/data/influencer-chat-mock";
import { initialInfluencerContracts } from "@/data/influencer-contracts-mock";
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
  const [invites, setInvites] = useState<InfluencerInvite[]>(initialInfluencerInvites);

  useEffect(() => {
    setInvites(loadInfluencerInvites());
  }, []);

  const pendingInvites = useMemo(() => invites.filter((i) => i.status === "pending"), [invites]);
  const dashboardInvites = useMemo(() => pendingInvites.slice(0, 3), [pendingInvites]);

  const pendingCount = invites.filter((i) => i.status === "pending").length;
  const activeCollabs = initialInfluencerContracts.filter((c) => c.status === "active").length;

  const setInviteStatus = (id: string, status: "accepted" | "rejected") => {
    setInvites((prev) => {
      const next = prev.map((i) => (i.id === id ? { ...i, status } : i));
      saveInfluencerInvites(next);
      return next;
    });
  };

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-white">Welcome back, Biruk</h1>
          <p className="mt-1 max-w-xl text-sm text-[#CBD5E1]">
            Manage your creator profile, campaign invites, contracts, and earnings.
          </p>
        </div>
        <Button
          asChild
          className="h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-5 text-white hover:bg-[color:var(--adfluence-orange-hover)]"
        >
          <Link href="/influencer/profile">Edit Profile</Link>
        </Button>
      </section>

      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <KpiCard label="Campaign Invites" value={String(pendingCount)} hint="Awaiting your response" />
        <KpiCard label="Active Collaborations" value={String(activeCollabs)} />
        <KpiCard label="Estimated Earnings" value="ETB 42K" hint="This quarter (mock)" />
        <KpiCard label="Profile Score" value="86%" hint="Strong visibility" />
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:p-6">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h2 className="text-lg font-semibold text-white">Profile completion</h2>
            <p className="mt-1 text-sm text-[#CBD5E1]">Complete your profile to unlock premium campaigns.</p>
          </div>
          <span className="rounded-full border border-white/15 bg-white/5 px-3 py-1 text-sm font-semibold text-white">
            86%
          </span>
        </div>
        <div className="mt-4 h-2 overflow-hidden rounded-full bg-white/10">
          <div className="h-full w-[86%] rounded-full bg-[color:var(--adfluence-orange)]" />
        </div>
        <ul className="mt-4 space-y-2 text-sm text-[#CBD5E1]">
          <li className="flex items-center gap-2">
            <span className="text-[color:var(--adfluence-orange)]">○</span> Add more content samples
          </li>
          <li className="flex items-center gap-2">
            <span className="text-[color:var(--adfluence-orange)]">○</span> Update pricing packages
          </li>
          <li className="flex items-center gap-2">
            <span className="text-[color:var(--adfluence-orange)]">○</span> Connect TikTok profile
          </li>
        </ul>
        <Button
          asChild
          className="mt-4 h-10 rounded-xl bg-white/10 text-white hover:bg-white/15"
        >
          <Link href="/influencer/profile">Complete Profile</Link>
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
          {dashboardInvites.length === 0 ? (
            <p className="text-sm text-[#CBD5E1]">No pending invites right now.</p>
          ) : (
            dashboardInvites.map((inv) => (
              <article
                key={inv.id}
                className="flex flex-col rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
              >
                <p className="font-semibold text-white">{inv.brand}</p>
                <p className="text-sm text-[#CBD5E1]">{inv.campaign}</p>
                <p className="mt-2 text-xs text-[#94A3B8]">Goal: {inv.goal}</p>
                <p className="text-xs text-[#94A3B8]">Budget: {inv.budget}</p>
                <p className="text-xs text-[#94A3B8]">Deadline: {inv.deadline}</p>
                <span className="mt-2 inline-flex w-fit rounded-full border border-amber-400/30 bg-amber-500/10 px-2 py-0.5 text-xs capitalize text-amber-200">
                  {inv.status}
                </span>
                <div className="mt-4 flex flex-wrap gap-2">
                  <Button
                    type="button"
                    size="sm"
                    className="h-9 flex-1 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                    onClick={() => setInviteStatus(inv.id, "accepted")}
                  >
                    Accept
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    className="h-9 flex-1 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                    onClick={() => setInviteStatus(inv.id, "rejected")}
                  >
                    Reject
                  </Button>
                  <Button
                    asChild
                    size="sm"
                    className="h-9 w-full rounded-xl bg-white/10 text-white hover:bg-white/15"
                  >
                    <Link href="/influencer/invites">View Details</Link>
                  </Button>
                </div>
              </article>
            ))
          )}
        </div>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:p-6">
        <h2 className="text-lg font-semibold text-white">Earnings overview</h2>
        <div className="mt-4 grid gap-3 sm:grid-cols-3">
          <div className="rounded-xl border border-white/10 bg-[#06111F] p-3">
            <p className="text-xs text-[#CBD5E1]">Pending payout</p>
            <p className="text-lg font-semibold text-amber-200">ETB 8,400</p>
          </div>
          <div className="rounded-xl border border-white/10 bg-[#06111F] p-3">
            <p className="text-xs text-[#CBD5E1]">Completed earnings</p>
            <p className="text-lg font-semibold text-emerald-300">ETB 34,200</p>
          </div>
          <div className="rounded-xl border border-white/10 bg-[#06111F] p-3">
            <p className="text-xs text-[#CBD5E1]">Next payout</p>
            <p className="text-lg font-semibold text-white">Jun 5, 2026</p>
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
          <ul className="mt-3 space-y-3">
            {initialCreatorConversations.slice(0, 3).map((c) => (
              <li
                key={c.id}
                className="flex items-start gap-3 rounded-xl border border-white/10 bg-[#06111F] p-3"
              >
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-[color:var(--adfluence-purple)] text-sm font-bold text-white">
                  {c.brandInitial}
                </div>
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-white">{c.businessName}</p>
                  <p className="truncate text-xs text-[#CBD5E1]">{c.lastMessage}</p>
                </div>
                {c.unreadCount > 0 ? (
                  <span className="shrink-0 rounded-full bg-[color:var(--adfluence-orange)] px-2 py-0.5 text-[10px] font-semibold">
                    {c.unreadCount}
                  </span>
                ) : null}
              </li>
            ))}
          </ul>
        </div>

        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <div className="flex items-center justify-between gap-2">
            <h2 className="text-lg font-semibold text-white">Active contracts</h2>
            <MessageSquare className="size-4 text-[color:var(--adfluence-purple-accent)]" aria-hidden />
          </div>
          <ul className="mt-3 space-y-2">
            {initialInfluencerContracts
              .filter((c) => c.status === "active")
              .map((c) => (
                <li
                  key={c.id}
                  className="flex items-center justify-between gap-2 rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm"
                >
                  <span className="truncate text-[#CBD5E1]">{c.brand}</span>
                  <span className="shrink-0 text-xs text-emerald-300">Active</span>
                </li>
              ))}
          </ul>
          <Button asChild className="mt-4 h-10 w-full rounded-xl bg-white/10 text-white hover:bg-white/15">
            <Link href="/influencer/contracts">View Contracts</Link>
          </Button>
        </div>
      </section>
    </div>
  );
}
