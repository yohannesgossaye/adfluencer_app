"use client";

import { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type AdminCampaignListItem,
  type AdminCampaignStats,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

const statuses = [
  { id: "draft", label: "Draft" },
  { id: "active", label: "Live" },
  { id: "paused", label: "Paused" },
  { id: "completed", label: "Completed" },
  { id: "cancelled", label: "Cancelled" },
];

const STATUS_LABEL: Record<string, string> = Object.fromEntries(statuses.map((s) => [s.id, s.label]));

function badge(s: string) {
  switch (s.toLowerCase()) {
    case "draft":
      return "bg-white/10 text-[#CBD5E1]";
    case "active":
      return "bg-emerald-500/15 text-emerald-200";
    case "completed":
      return "bg-sky-500/15 text-sky-200";
    case "paused":
      return "bg-rose-500/15 text-rose-200";
    default:
      return "bg-white/10 text-[#CBD5E1]";
  }
}

type Action = "pause" | "flag" | "resume" | "complete";

const ACTION_DONE: Record<Action, string> = {
  pause: "Campaign paused. Creators can no longer accept its invites.",
  flag: "Campaign flagged and paused.",
  resume: "Campaign is live again.",
  complete: "Campaign marked completed.",
};

export function AdminCampaigns() {
  const [rows, setRows] = useState<AdminCampaignListItem[]>([]);
  const [stats, setStats] = useState<AdminCampaignStats>({ active: 0, draft: 0, paused: 0 });
  const [status, setStatus] = useState<string>("all");
  const [niche, setNiche] = useState("");
  const [business, setBusiness] = useState("");
  const [filters, setFilters] = useState({ niche: "", business: "" });
  const [isLoading, setIsLoading] = useState(true);
  const [confirmComplete, setConfirmComplete] = useState<string | null>(null);
  const [busyId, setBusyId] = useState<string | null>(null);

  useEffect(() => {
    const t = setTimeout(() => setFilters({ niche: niche.trim(), business: business.trim() }), 300);
    return () => clearTimeout(t);
  }, [niche, business]);

  const loadCampaigns = useCallback(async () => {
    try {
      const res = await adminApi.listCampaigns({
        status: status === "all" ? undefined : status,
        niche: filters.niche || undefined,
        business: filters.business || undefined,
      });
      setRows(res.data?.items ?? []);
      if (res.data?.stats) setStats(res.data.stats);
    } catch {
      toast.error("Could not load campaigns.");
    } finally {
      setIsLoading(false);
    }
  }, [status, filters]);

  useEffect(() => {
    void Promise.resolve().then(loadCampaigns);
  }, [loadCampaigns]);

  const act = async (id: string, action: Action) => {
    setBusyId(id);
    try {
      if (action === "pause") await adminApi.pauseCampaign(id);
      else if (action === "flag") await adminApi.flagCampaign(id);
      else if (action === "resume") await adminApi.resumeCampaign(id);
      else await adminApi.completeCampaign(id);
      toast.success(ACTION_DONE[action]);
      setConfirmComplete(null);
      await loadCampaigns();
    } catch (err) {
      const msg = axios.isAxiosError(err) ? (err.response?.data as { message?: string } | undefined)?.message : undefined;
      toast.error(msg || "The action could not be completed.");
    } finally {
      setBusyId(null);
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
          { label: "Live", value: stats.active },
          { label: "Drafts", value: stats.draft },
          { label: "Paused", value: stats.paused },
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
            key={s.id}
            type="button"
            onClick={() => setStatus(s.id)}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm capitalize transition-colors sm:min-h-0 sm:px-3 sm:py-1.5 sm:text-xs",
              status === s.id
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                : "border-white/10 text-[#CBD5E1]",
            )}
          >
            {s.label}
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
                <span className={cn("rounded-full px-2 py-0.5 text-xs", badge(r.status))}>
                  {STATUS_LABEL[r.status] ?? r.status}
                </span>
              </div>
              <p className="mt-2 text-xs text-[#94A3B8]">
                {r.niche || "General"} · {r.budget_summary || "Flexible"} · {r.creators_invited} creators invited · Created {new Date(r.created_at).toLocaleDateString()}
              </p>
              <p className="mt-1 text-xs text-[#CBD5E1]">
                Contracts: {r.contract_summary} · Payment: {r.payment_summary}
              </p>
              {confirmComplete === r.id ? (
                <div className="mt-3 flex flex-wrap items-center gap-2 rounded-xl border border-white/15 bg-black/20 p-2">
                  <span className="text-xs text-[#E2E8F0]">Mark this campaign completed? This can’t be undone.</span>
                  <Button size="sm" disabled={busyId === r.id} className="h-11 bg-emerald-600 text-white hover:bg-emerald-700 sm:h-8" onClick={() => void act(r.id, "complete")}>
                    Mark completed
                  </Button>
                  <Button size="sm" variant="ghost" className="h-11 text-white hover:bg-white/10 sm:h-8" onClick={() => setConfirmComplete(null)}>
                    Cancel
                  </Button>
                </div>
              ) : r.status === "active" || r.status === "paused" ? (
                <div className="mt-3 flex flex-wrap gap-2">
                  {r.status === "active" ? (
                    <>
                      <Button size="sm" variant="outline" disabled={busyId === r.id} className="h-11 border-white/15 bg-transparent text-white sm:h-8" onClick={() => void act(r.id, "pause")}>
                        Pause
                      </Button>
                      <Button size="sm" disabled={busyId === r.id} className="h-11 bg-rose-500/20 text-rose-100 sm:h-8" onClick={() => void act(r.id, "flag")}>
                        Flag
                      </Button>
                    </>
                  ) : (
                    <Button size="sm" disabled={busyId === r.id} className="h-11 bg-[color:var(--adfluence-orange)] text-white sm:h-8" onClick={() => void act(r.id, "resume")}>
                      Resume
                    </Button>
                  )}
                  <Button size="sm" disabled={busyId === r.id} className="h-11 bg-emerald-500/20 text-emerald-100 sm:h-8" onClick={() => setConfirmComplete(r.id)}>
                    Mark completed
                  </Button>
                </div>
              ) : null}
            </article>
          ))
        )}
      </div>
    </div>
  );
}
