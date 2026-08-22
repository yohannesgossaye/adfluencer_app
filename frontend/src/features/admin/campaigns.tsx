"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type AdminCampaignListItem,
  type AdminCampaignStats,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

const statuses = [
  "draft",
  "pending",
  "active",
  "completed",
  "flagged",
];

function badge(s: string) {
  switch (s.toLowerCase()) {
    case "draft":
      return "bg-white/10 text-[#CBD5E1]";
    case "pending":
      return "bg-amber-500/15 text-amber-200";
    case "active":
      return "bg-emerald-500/15 text-emerald-200";
    case "completed":
      return "bg-sky-500/15 text-sky-200";
    case "flagged":
    case "paused":
      return "bg-rose-500/15 text-rose-200";
    default:
      return "bg-white/10";
  }
}

export function AdminCampaigns() {
  const [rows, setRows] = useState<AdminCampaignListItem[]>([]);
  const [stats, setStats] = useState<AdminCampaignStats>({ active: 0, pending: 0, flagged: 0 });
  const [status, setStatus] = useState<string>("all");
  const [niche, setNiche] = useState("");
  const [business, setBusiness] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const loadCampaigns = async () => {
    try {
      const res = await adminApi.listCampaigns({
        status: status === "all" ? undefined : status,
        niche: niche.trim() || undefined,
        business: business.trim() || undefined,
      });
      if (res.data?.items) {
        setRows(res.data.items);
      }
      if (res.data?.stats) {
        setStats(res.data.stats);
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadCampaigns();
  }, [status, niche, business]);

  const pauseCampaign = async (id: string) => {
    try {
      await adminApi.pauseCampaign(id);
      await loadCampaigns();
    } catch {
      // fallback
    }
  };

  const completeCampaign = async (id: string) => {
    try {
      await adminApi.completeCampaign(id);
      await loadCampaigns();
    } catch {
      // fallback
    }
  };

  const flagCampaign = async (id: string) => {
    try {
      await adminApi.flagCampaign(id);
      await loadCampaigns();
    } catch {
      // fallback
    }
  };

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Campaigns</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          View and moderate campaigns across all businesses on Adfluence.
        </p>
      </div>

      <section className="grid gap-3 sm:grid-cols-3">
        {[
          { label: "Active", value: stats.active },
          { label: "Pending", value: stats.pending },
          { label: "Flagged", value: stats.flagged },
        ].map((s) => (
          <article
            key={s.label}
            className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
          >
            <p className="text-xs text-[#CBD5E1]">{s.label}</p>
            <p className="mt-2 text-2xl font-bold text-[color:var(--adfluence-orange)]">{s.value}</p>
          </article>
        ))}
      </section>

      <div className="flex flex-wrap gap-2">
        <button
          type="button"
          onClick={() => setStatus("all")}
          className={cn(
            "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm transition-colors sm:min-h-0 sm:px-3 sm:py-1.5 sm:text-xs",
            status === "all"
              ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
              : "border-white/10 text-[#CBD5E1]",
          )}
        >
          All statuses
        </button>
        {statuses.map((s) => (
          <button
            key={s}
            type="button"
            onClick={() => setStatus(s)}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm capitalize transition-colors sm:min-h-0 sm:px-3 sm:py-1.5 sm:text-xs",
              status === s
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                : "border-white/10 text-[#CBD5E1]",
            )}
          >
            {s}
          </button>
        ))}
      </div>

      <div className="flex flex-wrap gap-3">
        <Input
          value={niche}
          onChange={(e) => setNiche(e.target.value)}
          placeholder="Filter niche"
          className="h-11 max-w-xs rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white sm:h-10"
        />
        <Input
          value={business}
          onChange={(e) => setBusiness(e.target.value)}
          placeholder="Filter business"
          className="h-11 max-w-xs rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white sm:h-10"
        />
      </div>

      <div className="space-y-3 min-w-0">
        {rows.length === 0 && !isLoading ? (
          <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center">
            <p className="text-white">No campaigns found.</p>
          </div>
        ) : (
          rows.map((r) => (
            <article
              key={r.id}
              className="w-full min-w-0 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
            >
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <p className="font-semibold text-white">{r.title}</p>
                  <p className="text-sm text-[#CBD5E1]">{r.business_name}</p>
                </div>
                <span className={cn("rounded-full px-2 py-0.5 text-xs capitalize", badge(r.status))}>
                  {r.status}
                </span>
              </div>
              <p className="mt-2 text-xs text-[#94A3B8]">
                {r.niche || "General"} · {r.budget_summary || "Flexible"} · {r.creators_invited} creators invited · Created {new Date(r.created_at).toLocaleDateString()}
              </p>
              <p className="mt-1 text-xs text-[#CBD5E1]">
                Contracts: {r.contract_summary} · Payment: {r.payment_summary}
              </p>
              <div className="mt-3 flex flex-wrap gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  className="h-11 sm:h-8 border-white/15 text-white"
                  disabled={r.status === "completed"}
                  onClick={() => pauseCampaign(r.id)}
                >
                  Pause
                </Button>
                <Button
                  size="sm"
                  className="h-11 sm:h-8 bg-emerald-500/20 text-emerald-100"
                  disabled={r.status === "completed"}
                  onClick={() => completeCampaign(r.id)}
                >
                  Mark completed
                </Button>
                <Button
                  size="sm"
                  className="h-11 sm:h-8 bg-rose-500/20 text-rose-100"
                  disabled={r.status === "flagged"}
                  onClick={() => flagCampaign(r.id)}
                >
                  Flag
                </Button>
              </div>
            </article>
          ))
        )}
      </div>
    </div>
  );
}
