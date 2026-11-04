"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import { ListStatus } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import type { LoadState } from "@/hooks/use-paged-list";
import {
  adminApi,
  type AdminDashboardSummaryResponse,
} from "@/lib/api/admin";
import { adminErrorMessage } from "@/lib/api-errors";
import { formatEtb } from "@/lib/contracts";
import { cn } from "@/lib/utils";

/** One "needs action" tile: amber when there is work waiting, calm when clear. */
function ActionTile({ label, count, href, hint }: { label: string; count: number; href: string; hint: string }) {
  const waiting = count > 0;
  return (
    <Link
      href={href}
      className={cn(
        "group flex flex-col justify-between rounded-2xl border p-4 transition-colors focus-visible:outline-2 focus-visible:outline-[#A78BFA]",
        waiting
          ? "border-amber-300/40 bg-amber-300/10 hover:bg-amber-300/15"
          : "border-white/10 bg-[rgba(255,255,255,0.04)] hover:bg-white/[0.07]",
      )}
    >
      <p className="text-sm font-medium text-[#CBD5E1]">{label}</p>
      <p className={cn("mt-2 text-3xl font-bold tabular-nums", waiting ? "text-amber-200" : "text-white/60")}>{count}</p>
      <p className="mt-1 text-xs text-[#94A3B8] group-hover:text-white">{waiting ? `${hint} →` : "All clear"}</p>
    </Link>
  );
}

function Kpi({ label, value }: { label: string; value: string }) {
  return (
    <article className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
      <p className="text-[11px] font-medium uppercase leading-tight tracking-wide text-[#94A3B8]">
        {label}
      </p>
      <p className="mt-2 text-xl font-bold text-[color:var(--adfluence-orange)] md:text-2xl">
        {value}
      </p>
    </article>
  );
}

export function AdminDashboard() {
  const pathname = usePathname();
  const [summary, setSummary] = useState<AdminDashboardSummaryResponse | null>(null);
  const [state, setState] = useState<LoadState>("loading");

  const loadSummary = useCallback(async () => {
    try {
      const res = await adminApi.getDashboardSummary();
      setSummary(res.data ?? null);
      setState("ready");
    } catch {
      setState("error");
    }
  }, []);

  useEffect(() => {
    void Promise.resolve().then(loadSummary);
  }, [pathname, loadSummary]);

  const approve = async (userId: string) => {
    try {
      await adminApi.approveInfluencerApplication(userId);
      toast.success("Creator approved.");
      await loadSummary();
    } catch (err) {
      toast.error(adminErrorMessage(err, "Could not approve this creator."));
    }
  };

  const kpis = useMemo(() => {
    const k = summary?.kpis;
    // Totals only; anything that needs an admin's action is in the row above.
    return [
      { label: "Creators", value: String(k?.total_influencers ?? 0) },
      { label: "Businesses", value: String(k?.total_businesses ?? 0) },
      { label: "Live campaigns", value: String(k?.active_campaigns ?? 0) },
      { label: "All campaigns", value: String(k?.total_campaigns ?? 0) },
      { label: "Completed contracts", value: String(k?.completed_contracts ?? 0) },
      { label: "Revenue this month", value: formatEtb(Math.round(k?.monthly_platform_revenue_etb ?? 0)) },
      { label: "Blocked accounts", value: String(k?.blocked_users ?? 0) },
      { label: "All accounts", value: String(k?.total_users ?? 0) },
    ];
  }, [summary]);

  const previewPending = summary?.pending_approvals || [];
  // Logins are noise here; show actions people took.
  const activity = (summary?.recent_activity || []).filter((a) => !/^log(in|out)$/i.test(a.label.trim()));
  const alerts = summary?.alerts || [];

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Admin Dashboard</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Monitor platform activity, approvals, campaigns, and marketplace health.
        </p>
      </div>

      <ListStatus
        state={state}
        isEmpty={false}
        emptyText=""
        errorText="The dashboard could not be loaded."
        onRetry={() => {
          setState("loading");
          void loadSummary();
        }}
      >
        <section aria-labelledby="needs-action">
          <h2 id="needs-action" className="text-lg font-semibold text-white">
            Needs action today
          </h2>
          <div className="mt-3 grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
            <ActionTile label="Payments to confirm" count={summary?.kpis?.payments_to_confirm ?? 0} href="/admin/payments" hint="Confirm payments" />
            <ActionTile label="Businesses to verify" count={summary?.kpis?.businesses_to_review ?? 0} href="/admin/businesses" hint="Review businesses" />
            <ActionTile label="Creators to approve" count={summary?.kpis?.pending_influencer_approvals ?? 0} href="/admin/influencers" hint="Review creators" />
            <ActionTile label="Open disputes" count={summary?.kpis?.open_disputes ?? 0} href="/admin/disputes" hint="Handle disputes" />
            <ActionTile label="Urgent tickets" count={summary?.kpis?.urgent_support_tickets ?? 0} href="/admin/support" hint="Answer tickets" />
          </div>
        </section>

        <h2 className="sr-only">Platform totals</h2>
        <section className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          {kpis.map((k) => (
            <Kpi key={k.label} label={k.label} value={k.value} />
          ))}
        </section>

        <div className="grid gap-6 xl:grid-cols-2">
          <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
            <h2 className="text-lg font-semibold text-white">Creators waiting for approval</h2>
            {previewPending.length === 0 ? (
              <p className="mt-3 text-sm text-[#CBD5E1]">No pending approvals.</p>
            ) : (
              <ul className="mt-3 space-y-3">
                {previewPending.map((a) => (
                  <li
                    key={a.user_id}
                    className="flex flex-col justify-between gap-2 rounded-xl border border-white/10 bg-[#06111F] p-3 sm:flex-row sm:items-center"
                  >
                    <div>
                      <p className="font-medium text-white">{a.name}</p>
                      <p className="text-xs text-[#CBD5E1]">
                        {a.niche || "General"} · {a.platform || "Multi"} · {a.followers.toLocaleString()} followers
                      </p>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        type="button"
                        size="sm"
                        className="h-11 rounded-lg bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] sm:h-8"
                        onClick={() => void approve(a.user_id)}
                      >
                        Approve
                      </Button>
                      <Button
                        asChild
                        size="sm"
                        variant="outline"
                        className="h-11 rounded-lg border-white/15 bg-transparent text-white hover:bg-white/10 sm:h-8"
                      >
                        <Link href="/admin/influencers">Review</Link>
                      </Button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </section>

          <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
            <h2 className="text-lg font-semibold text-white">Platform activity</h2>
            {activity.length === 0 ? (
              <p className="mt-3 text-sm text-[#CBD5E1]">No recent actions yet. Approvals, payments and moderation will appear here.</p>
            ) : (
              <ul className="mt-3 space-y-2">
                {activity.map((item) => (
                  <li
                    key={item.id}
                    className="flex items-start justify-between gap-2 border-b border-white/5 py-2 text-sm last:border-0"
                  >
                    <div>
                      <p className="font-medium text-white">{item.label}</p>
                      <p className="text-xs text-[#CBD5E1]">{item.detail}</p>
                    </div>
                    <span className="shrink-0 text-[10px] text-[#94A3B8]">
                      {new Date(item.created_at).toLocaleDateString()}
                    </span>
                  </li>
                ))}
              </ul>
            )}
          </section>
        </div>

        <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:p-6">
          <h2 className="text-lg font-semibold text-white">Revenue / escrow performance</h2>
          <div className="mt-4 grid gap-3 sm:grid-cols-2">
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs text-[#CBD5E1]">Total campaign escrow value</p>
              <p className="text-lg font-semibold text-white">
                ETB {Math.round(summary?.revenue_summary?.total_campaign_value_etb ?? 0).toLocaleString()}
              </p>
            </div>
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs text-[#CBD5E1]">Total platform commission (5% + 5%)</p>
              <p className="text-lg font-semibold text-[color:var(--adfluence-orange)]">
                ETB {Math.round(summary?.revenue_summary?.total_platform_fees_etb ?? 0).toLocaleString()}
              </p>
            </div>
          </div>
        </section>

        {alerts.length > 0 ? (
          <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
            <h2 className="text-lg font-semibold text-white">Alerts</h2>
            <ul className="mt-3 space-y-2">
              {alerts.map((al) => (
                <li
                  key={al.id}
                  className={cn(
                    "rounded-xl border px-3 py-2 text-sm",
                    al.severity === "high"
                      ? "border-rose-400/30 bg-rose-500/10 text-rose-100"
                      : al.severity === "medium"
                        ? "border-amber-400/25 bg-amber-500/10 text-amber-100"
                        : "border-white/10 bg-[color:var(--adfluence-surface)] text-[#CBD5E1]",
                  )}
                >
                  <span className="font-medium text-white">{al.title}</span> — {al.detail}
                </li>
              ))}
            </ul>
          </section>
        ) : null}
      </ListStatus>
    </div>
  );
}
