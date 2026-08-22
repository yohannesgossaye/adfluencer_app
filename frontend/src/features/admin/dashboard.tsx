"use client";

import { usePathname } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  adminApi,
  type AdminDashboardSummaryResponse,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

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
  const [isLoading, setIsLoading] = useState(true);

  const loadSummary = async () => {
    try {
      const res = await adminApi.getDashboardSummary();
      if (res.data) {
        setSummary(res.data);
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadSummary();
  }, [pathname]);

  const setAppStatus = async (userId: string, status: "approved" | "rejected") => {
    try {
      if (status === "approved") {
        await adminApi.approveInfluencerApplication(userId);
      } else {
        await adminApi.rejectInfluencerApplication(userId, { reason: "Application rejected by admin." });
      }
      await loadSummary();
    } catch {
      // fallback
    }
  };

  const kpis = useMemo(() => {
    const k = summary?.kpis;
    return [
      { label: "Total Influencers", value: String(k?.total_influencers ?? 0) },
      { label: "Total Businesses", value: String(k?.total_businesses ?? 0) },
      { label: "Pending Approvals", value: String(k?.pending_influencer_approvals ?? 0) },
      { label: "Active Campaigns", value: String(k?.active_campaigns ?? 0) },
      { label: "Total Campaigns", value: String(k?.total_campaigns ?? 0) },
      { label: "Total Users", value: String(k?.total_users ?? 0) },
      { label: "Open Disputes", value: String(k?.open_disputes ?? 0) },
      { label: "Support Tickets", value: String(k?.open_support_tickets ?? 0) },
      {
        label: "Monthly Revenue",
        value: `ETB ${Math.round(k?.monthly_platform_revenue_etb ?? 0).toLocaleString()}`,
      },
      { label: "Platform Fee", value: "5% + 5%" },
      { label: "Blocked Users", value: String(k?.blocked_users ?? 0) },
      { label: "Completed Contracts", value: String(k?.completed_contracts ?? 0) },
    ];
  }, [summary]);

  const previewPending = summary?.pending_approvals || [];
  const activity = summary?.recent_activity || [];
  const alerts = summary?.alerts || [];

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Admin Dashboard</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Monitor platform activity, approvals, campaigns, and marketplace health.
        </p>
      </div>

      <section className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {kpis.map((k) => (
          <Kpi key={k.label} label={k.label} value={k.value} />
        ))}
      </section>

      <div className="grid gap-6 xl:grid-cols-2">
        <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <h2 className="text-lg font-semibold text-white">Pending influencer approvals</h2>
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
                      className="h-11 rounded-lg bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:h-8"
                      onClick={() => setAppStatus(a.user_id, "approved")}
                    >
                      Approve
                    </Button>
                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      className="h-11 rounded-lg border-white/15 text-white hover:bg-white/10 sm:h-8"
                      onClick={() => setAppStatus(a.user_id, "rejected")}
                    >
                      Reject
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
            <p className="mt-3 text-sm text-[#CBD5E1]">No recent activity logs.</p>
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
    </div>
  );
}
