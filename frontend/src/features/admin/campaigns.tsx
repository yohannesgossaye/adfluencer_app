"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { ConfirmActionDialog, type ConfirmActionRequest } from "@/components/common/confirm-action-dialog";
import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type AdminCampaignListItem,
  type AdminCampaignStats,
} from "@/lib/api/admin";
import { usePagedList } from "@/hooks/use-paged-list";
import { adminErrorMessage } from "@/lib/api-errors";
import { cn } from "@/lib/utils";
import { formatDate } from "@/lib/contracts";
import { StatusBadge } from "@/components/common/status-badge";
import { campaignStatus } from "@/lib/status";

const statuses = [
  { id: "draft", label: "Draft" },
  { id: "active", label: "Live" },
  { id: "paused", label: "Paused" },
  { id: "completed", label: "Completed" },
  { id: "cancelled", label: "Cancelled" },
];



export function AdminCampaigns() {
  const [status, setStatus] = useState<string>("all");
  const [niche, setNiche] = useState("");
  const [business, setBusiness] = useState("");
  const [filters, setFilters] = useState({ niche: "", business: "" });
  const [confirm, setConfirm] = useState<ConfirmActionRequest | null>(null);
  const [busyId, setBusyId] = useState<string | null>(null);

  useEffect(() => {
    const t = setTimeout(() => setFilters({ niche: niche.trim(), business: business.trim() }), 300);
    return () => clearTimeout(t);
  }, [niche, business]);

  const fetchPage = useCallback(
    async (page: number) =>
      (
        await adminApi.listCampaigns({
          status: status === "all" ? undefined : status,
          niche: filters.niche || undefined,
          business: filters.business || undefined,
          page,
        })
      ).data,
    [status, filters],
  );
  const list = usePagedList(fetchPage);
  const rows: AdminCampaignListItem[] = list.items;
  const stats: AdminCampaignStats = list.last?.stats ?? { active: 0, draft: 0, paused: 0 };
  const loadCampaigns = list.reload;

  const run = async (id: string, action: () => Promise<unknown>, success: string) => {
    setBusyId(id);
    try {
      await action();
      toast.success(success);
      await loadCampaigns();
    } catch (err) {
      toast.error(adminErrorMessage(err, "The action could not be completed."));
      throw err;
    } finally {
      setBusyId(null);
    }
  };

  const askPause = (r: AdminCampaignListItem) =>
    setConfirm({
      title: `Pause "${r.title}"?`,
      description: "Creators can no longer accept its invites until you resume it. Existing contracts are not changed.",
      confirmLabel: "Pause campaign",
      tone: "danger",
      reasonLabel: "Reason (kept in the audit log)",
      reasonPlaceholder: "e.g. Misleading product claims, reported by a creator",
      onConfirm: (reason) =>
        run(r.id, () => adminApi.pauseCampaign(r.id, reason), "Campaign paused. Creators can no longer accept its invites."),
    });

  const askComplete = (r: AdminCampaignListItem) =>
    setConfirm({
      title: `Mark "${r.title}" completed?`,
      description: "This closes the campaign for good and can't be undone.",
      confirmLabel: "Mark completed",
      reasonLabel: "Reason (kept in the audit log)",
      reasonPlaceholder: "e.g. Business asked us to close it",
      onConfirm: (reason) => run(r.id, () => adminApi.completeCampaign(r.id, reason), "Campaign marked completed."),
    });

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
          { label: "Live", value: list.state === "ready" ? stats.active : "—" },
          { label: "Drafts", value: list.state === "ready" ? stats.draft : "—" },
          { label: "Paused", value: list.state === "ready" ? stats.paused : "—" },
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
              ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
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
                ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
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
          aria-label="Filter by niche"
          className="h-11 max-w-xs rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white sm:h-10"
        />
        <Input
          value={business}
          onChange={(e) => setBusiness(e.target.value)}
          placeholder="Filter business"
          aria-label="Filter by business"
          className="h-11 max-w-xs rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white sm:h-10"
        />
      </div>

      <div className="space-y-3 min-w-0">
        <ListStatus
          state={list.state}
          isEmpty={rows.length === 0}
          emptyText="No campaigns match these filters."
          errorText="Campaigns could not be loaded."
          onRetry={list.retry}
        >
          {rows.map((r) => (
            <article
              key={r.id}
              className="w-full min-w-0 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
            >
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <p className="font-semibold text-white">{r.title}</p>
                  <p className="text-sm text-[#CBD5E1]">{r.business_name}</p>
                </div>
                <StatusBadge status={campaignStatus(r.status)} />
              </div>
              <p className="mt-2 text-xs text-[#94A3B8]">
                {r.niche || "General"} · {r.budget_summary || "Flexible"} · {r.creators_invited} creators invited · Created {formatDate(r.created_at)}
              </p>
              <p className="mt-1 text-xs text-[#CBD5E1]">
                Contracts: {r.contract_summary} · Payment: {r.payment_summary}
              </p>
              {r.status === "active" || r.status === "paused" ? (
                <div className="mt-3 flex flex-wrap gap-2">
                  {r.status === "active" ? (
                    <Button size="sm" disabled={busyId === r.id} className="h-11 border border-rose-400/50 bg-transparent text-rose-200 hover:bg-rose-500/10 sm:h-8" onClick={() => askPause(r)}>
                      Pause…
                    </Button>
                  ) : (
                    <Button
                      size="sm"
                      disabled={busyId === r.id}
                      className="h-11 bg-[color:var(--adfluence-action)] text-white sm:h-8"
                      onClick={() => void run(r.id, () => adminApi.resumeCampaign(r.id), "Campaign is live again.").catch(() => {})}
                    >
                      Resume
                    </Button>
                  )}
                  <Button size="sm" variant="outline" disabled={busyId === r.id} className="h-11 border-white/15 bg-transparent text-white sm:h-8" onClick={() => askComplete(r)}>
                    Mark completed…
                  </Button>
                </div>
              ) : null}
            </article>
          ))}
          <LoadMore
            shown={rows.length}
            total={list.total}
            hasMore={list.hasMore}
            loading={list.loadingMore}
            onLoadMore={list.loadMore}
          />
        </ListStatus>
      </div>

      <ConfirmActionDialog request={confirm} onClose={() => setConfirm(null)} />
    </div>
  );
}
