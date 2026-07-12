"use client";

import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  loadAdminCampaigns,
  saveAdminCampaigns,
  type AdminCampaignRecord,
  type AdminCampaignStatus,
} from "@/data/admin-campaigns-mock";
import { cn } from "@/lib/utils";

const statuses: AdminCampaignStatus[] = [
  "draft",
  "pending",
  "active",
  "completed",
  "flagged",
];

function badge(s: AdminCampaignStatus) {
  switch (s) {
    case "draft":
      return "bg-white/10 text-[#CBD5E1]";
    case "pending":
      return "bg-amber-500/15 text-amber-200";
    case "active":
      return "bg-emerald-500/15 text-emerald-200";
    case "completed":
      return "bg-sky-500/15 text-sky-200";
    case "flagged":
      return "bg-rose-500/15 text-rose-200";
    default:
      return "bg-white/10";
  }
}

export function AdminCampaigns() {
  const [rows, setRows] = useState<AdminCampaignRecord[]>([]);
  const [status, setStatus] = useState<"all" | AdminCampaignStatus>("all");
  const [niche, setNiche] = useState("");
  const [business, setBusiness] = useState("");

  useEffect(() => {
    setRows(loadAdminCampaigns());
  }, []);

  const filtered = useMemo(() => {
    const n = niche.trim().toLowerCase();
    const b = business.trim().toLowerCase();
    return rows.filter((r) => {
      const st = status === "all" || r.status === status;
      const ni = !n || r.niche.toLowerCase().includes(n);
      const bu = !b || r.business.toLowerCase().includes(b);
      return st && ni && bu;
    });
  }, [rows, status, niche, business]);

  const persist = (next: AdminCampaignRecord[]) => {
    saveAdminCampaigns(next);
    setRows(next);
  };

  const setRowStatus = (id: string, nextStatus: AdminCampaignStatus) => {
    persist(rows.map((r) => (r.id === id ? { ...r, status: nextStatus } : r)));
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
          { label: "Active", value: rows.filter((r) => r.status === "active").length },
          { label: "Pending", value: rows.filter((r) => r.status === "pending").length },
          { label: "Flagged", value: rows.filter((r) => r.status === "flagged").length },
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
            "rounded-full border px-3 py-1.5 text-xs",
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
              "rounded-full border px-3 py-1.5 text-xs capitalize",
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
          className="h-10 max-w-xs rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white"
        />
        <Input
          value={business}
          onChange={(e) => setBusiness(e.target.value)}
          placeholder="Filter business"
          className="h-10 max-w-xs rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white"
        />
      </div>

      <div className="space-y-3 overflow-x-auto">
        {filtered.map((r) => (
          <article
            key={r.id}
            className="min-w-[300px] rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
          >
            <div className="flex flex-wrap items-start justify-between gap-2">
              <div>
                <p className="font-semibold text-white">{r.title}</p>
                <p className="text-sm text-[#CBD5E1]">{r.business}</p>
              </div>
              <span className={cn("rounded-full px-2 py-0.5 text-xs capitalize", badge(r.status))}>
                {r.status}
              </span>
            </div>
            <p className="mt-2 text-xs text-[#94A3B8]">
              {r.niche} · {r.budget} · {r.creatorsInvited} creators invited · Created {r.createdAt}
            </p>
            <p className="mt-1 text-xs text-[#CBD5E1]">
              Contracts: {r.contractStatus} · Payment: {r.paymentStatus}
            </p>
            <div className="mt-3 flex flex-wrap gap-2">
              <Button size="sm" className="h-8 bg-white/10 text-white">
                View
              </Button>
              <Button
                size="sm"
                variant="outline"
                className="h-8 border-white/15 text-white"
                disabled={r.status === "completed"}
                onClick={() => setRowStatus(r.id, "pending")}
              >
                Pause
              </Button>
              <Button
                size="sm"
                className="h-8 bg-emerald-500/20 text-emerald-100"
                disabled={r.status === "completed"}
                onClick={() => setRowStatus(r.id, "completed")}
              >
                Mark completed
              </Button>
              <Button
                size="sm"
                className="h-8 bg-rose-500/20 text-rose-100"
                disabled={r.status === "flagged"}
                onClick={() => setRowStatus(r.id, "flagged")}
              >
                Flag
              </Button>
            </div>
          </article>
        ))}
      </div>
    </div>
  );
}
