"use client";

import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  initialInfluencerInvites,
  loadInfluencerInvites,
  saveInfluencerInvites,
  type InfluencerInvite,
  type InfluencerInviteStatus,
} from "@/data/influencer-invites-mock";
import { cn } from "@/lib/utils";

const tabs: Array<{ id: "all" | InfluencerInviteStatus; label: string }> = [
  { id: "all", label: "All" },
  { id: "pending", label: "Pending" },
  { id: "accepted", label: "Accepted" },
  { id: "rejected", label: "Rejected" },
];

function statusStyle(status: InfluencerInviteStatus) {
  switch (status) {
    case "pending":
      return "border-amber-400/30 bg-amber-500/10 text-amber-200";
    case "accepted":
      return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
    case "rejected":
      return "border-rose-400/30 bg-rose-500/10 text-rose-200";
    default:
      return "border-white/10 bg-white/5 text-[#CBD5E1]";
  }
}

export function InfluencerInvites() {
  const [invites, setInvites] = useState<InfluencerInvite[]>(initialInfluencerInvites);
  const [tab, setTab] = useState<(typeof tabs)[number]["id"]>("all");

  useEffect(() => {
    setInvites(loadInfluencerInvites());
  }, []);

  const filtered = useMemo(() => {
    if (tab === "all") return invites;
    return invites.filter((i) => i.status === tab);
  }, [invites, tab]);

  const updateStatus = (id: string, status: InfluencerInviteStatus) => {
    setInvites((prev) => {
      const next = prev.map((i) => (i.id === id ? { ...i, status } : i));
      saveInfluencerInvites(next);
      return next;
    });
  };

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Campaign invites</h1>
        <p className="mt-1 text-sm text-[#CBD5E1]">
          Review invitations from brands and respond on your timeline.
        </p>
      </div>

      <div className="flex flex-wrap gap-2">
        {tabs.map((t) => (
          <button
            key={t.id}
            type="button"
            onClick={() => setTab(t.id)}
            className={cn(
              "rounded-full border px-4 py-2 text-sm font-medium transition-colors",
              tab === t.id
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                : "border-white/10 bg-[rgba(255,255,255,0.06)] text-[#CBD5E1] hover:border-white/20 hover:text-white",
            )}
          >
            {t.label}
          </button>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {filtered.map((inv) => (
          <article
            key={inv.id}
            className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5"
          >
            <div className="flex flex-wrap items-start justify-between gap-2">
              <div>
                <p className="font-semibold text-white">{inv.brand}</p>
                <p className="text-sm text-[#CBD5E1]">{inv.campaign}</p>
              </div>
              <span
                className={cn("rounded-full border px-2.5 py-0.5 text-xs capitalize", statusStyle(inv.status))}
              >
                {inv.status}
              </span>
            </div>
            <dl className="mt-3 space-y-1 text-sm text-[#CBD5E1]">
              <div>
                <dt className="inline text-[#94A3B8]">Goal: </dt>
                <dd className="inline text-white">{inv.goal}</dd>
              </div>
              <div>
                <dt className="inline text-[#94A3B8]">Budget: </dt>
                <dd className="inline text-white">{inv.budget}</dd>
              </div>
              <div>
                <dt className="inline text-[#94A3B8]">Deliverables: </dt>
                <dd className="inline text-white">{inv.deliverables}</dd>
              </div>
              <div>
                <dt className="inline text-[#94A3B8]">Timeline: </dt>
                <dd className="inline text-white">{inv.timeline}</dd>
              </div>
            </dl>
            <p className="mt-3 rounded-xl border border-white/10 bg-[#06111F] p-3 text-sm text-[#CBD5E1]">
              {inv.messagePreview}
            </p>
            {inv.status === "pending" ? (
              <div className="mt-4 flex flex-col gap-2 sm:flex-row">
                <Button
                  type="button"
                  className="h-10 flex-1 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                  onClick={() => updateStatus(inv.id, "accepted")}
                >
                  Accept
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  className="h-10 flex-1 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                  onClick={() => updateStatus(inv.id, "rejected")}
                >
                  Reject
                </Button>
              </div>
            ) : null}
          </article>
        ))}
      </div>
    </div>
  );
}
