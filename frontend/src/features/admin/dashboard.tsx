"use client";

import { usePathname } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { adminActivityFeed, adminAlerts } from "@/data/admin-activity-mock";
import { adminPlatformKpis } from "@/data/admin-kpis";
import {
  initialAdminApplications,
  loadApplications,
  saveApplications,
  type AdminInfluencerApplication,
} from "@/data/admin-influencer-applications";
import { initialAdminDisputes, loadDisputes, type AdminDispute } from "@/data/admin-disputes-mock";
import {
  initialSupportTickets,
  loadSupportTickets,
  type SupportTicket,
} from "@/data/support";
import { getBlockedEmails } from "@/lib/blocked-users";
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
  const [applications, setApplications] =
    useState<AdminInfluencerApplication[]>(initialAdminApplications);
  const [disputes, setDisputes] = useState<AdminDispute[]>(initialAdminDisputes);
  const [tickets, setTickets] = useState<SupportTicket[]>(initialSupportTickets);
  const [blockedCount, setBlockedCount] = useState(0);

  useEffect(() => {
    const refresh = () => {
      setApplications(loadApplications());
      setDisputes(loadDisputes());
      setTickets(loadSupportTickets());
      setBlockedCount(getBlockedEmails().length);
    };
    refresh();
  }, [pathname]);

  const pendingApprovals = useMemo(
    () => applications.filter((a) => a.status === "pending").length,
    [applications],
  );
  const openDisputes = useMemo(
    () => disputes.filter((d) => d.status === "open" || d.status === "under_review").length,
    [disputes],
  );
  const openTickets = useMemo(
    () =>
      tickets.filter(
        (t) => t.status === "open" || t.status === "in_review" || t.status === "escalated",
      ).length,
    [tickets],
  );

  const previewPending = useMemo(
    () => applications.filter((a) => a.status === "pending").slice(0, 3),
    [applications],
  );

  const setAppStatus = (id: string, status: "approved" | "rejected") => {
    setApplications((prev) => {
      const next = prev.map((a) => (a.id === id ? { ...a, status } : a));
      saveApplications(next);
      return next;
    });
  };

  const kpis = useMemo(
    () => [
      { label: "Total Influencers", value: String(adminPlatformKpis.totalInfluencers) },
      { label: "Total Businesses", value: String(adminPlatformKpis.totalBusinesses) },
      { label: "Pending Influencer Approvals", value: String(pendingApprovals) },
      { label: "Active Campaigns", value: String(adminPlatformKpis.activeCampaigns) },
      { label: "Total Campaigns", value: String(adminPlatformKpis.totalCampaigns) },
      { label: "Total Users", value: String(adminPlatformKpis.totalUsers) },
      { label: "Open Disputes", value: String(openDisputes) },
      { label: "Support Tickets", value: String(openTickets) },
      { label: "Monthly Revenue", value: adminPlatformKpis.monthlyRevenue },
      { label: "Platform Satisfaction", value: adminPlatformKpis.platformSatisfaction },
      { label: "Blocked Users", value: String(blockedCount) },
      { label: "Completed Contracts", value: adminPlatformKpis.completedContracts.toLocaleString() },
    ],
    [pendingApprovals, openDisputes, openTickets, blockedCount],
  );

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
          <h2 className="text-lg font-semibold text-white">Pending approvals</h2>
          <p className="mt-1 text-xs text-[#CBD5E1]">Influencer applications awaiting review</p>
          <ul className="mt-4 space-y-3">
            {previewPending.map((a) => (
              <li
                key={a.id}
                className="flex flex-col gap-2 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3 sm:flex-row sm:items-center sm:justify-between"
              >
                <div>
                  <p className="font-medium text-white">{a.name}</p>
                  <p className="text-xs text-[#CBD5E1]">
                    {a.niche} · {a.platform} · {a.followers.toLocaleString()} followers
                  </p>
                </div>
                <div className="flex gap-2">
                  <Button
                    type="button"
                    size="sm"
                    className="h-11 rounded-lg bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:h-8"
                    onClick={() => setAppStatus(a.id, "approved")}
                  >
                    Approve
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    className="h-11 rounded-lg border-white/15 text-white hover:bg-white/10 sm:h-8"
                    onClick={() => setAppStatus(a.id, "rejected")}
                  >
                    Reject
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </section>

        <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <h2 className="text-lg font-semibold text-white">Platform activity</h2>
          <ul className="mt-3 space-y-2">
            {adminActivityFeed.map((item) => (
              <li
                key={item.id}
                className="flex items-start justify-between gap-2 border-b border-white/5 py-2 text-sm last:border-0"
              >
                <div>
                  <p className="font-medium text-white">{item.label}</p>
                  <p className="text-xs text-[#CBD5E1]">{item.detail}</p>
                </div>
                <span className="shrink-0 text-[10px] text-[#94A3B8]">{item.time}</span>
              </li>
            ))}
          </ul>
        </section>
      </div>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:p-6">
        <h2 className="text-lg font-semibold text-white">Revenue / campaign performance</h2>
        <div className="mt-4 grid gap-3 sm:grid-cols-2">
          <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
            <p className="text-xs text-[#CBD5E1]">Total campaign value (mock)</p>
            <p className="text-lg font-semibold text-white">ETB 4.2M</p>
          </div>
          <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
            <p className="text-xs text-[#CBD5E1]">Platform commission (placeholder)</p>
            <p className="text-lg font-semibold text-[color:var(--adfluence-orange)]">12.5%</p>
          </div>
        </div>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
        <h2 className="text-lg font-semibold text-white">Alerts</h2>
        <ul className="mt-3 space-y-2">
          {adminAlerts.map((al) => (
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
    </div>
  );
}
