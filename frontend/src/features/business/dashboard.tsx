"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import {
  ArrowRight,
  BarChart3,
  FileText,
  MessageSquare,
  Sparkles,
  Target,
  Users,
} from "lucide-react";

import { VerificationStatusCard } from "@/components/business/verification-status-card";
import { RecommendedInfluencerCard } from "@/components/cards/recommended-influencer-card";
import { GradientButton } from "@/components/common/gradient-button";
import { ListStatus } from "@/components/common/list-status";
import { BusinessKpiCard } from "@/components/dashboard/business-kpi-card";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { LoadState } from "@/hooks/use-paged-list";
import { businessApi, type DashboardSummaryResponse } from "@/lib/api/business";
import { chatApi, type ChatConversationItem } from "@/lib/api/chat";
import { formatDate } from "@/lib/contracts";
import { StatusBadge } from "@/components/common/status-badge";
import { campaignStatus } from "@/lib/status";


export function BusinessDashboard() {
  const [summary, setSummary] = useState<DashboardSummaryResponse | null>(null);
  const [conversations, setConversations] = useState<ChatConversationItem[]>([]);
  const [state, setState] = useState<LoadState>("loading");

  const loadDashboard = useCallback(async () => {
    const [sumRes, chatRes] = await Promise.allSettled([
      businessApi.getDashboardSummary(),
      chatApi.listConversations(),
    ]);
    if (chatRes.status === "fulfilled" && chatRes.value.data) {
      setConversations(Array.isArray(chatRes.value.data) ? chatRes.value.data : []);
    }
    // The summary drives every number here; without it, say so instead of showing zeros.
    if (sumRes.status === "fulfilled") {
      setSummary(sumRes.value.data ?? null);
      setState("ready");
    } else {
      setState("error");
    }
  }, []);

  useEffect(() => {
    void Promise.resolve().then(loadDashboard);
  }, [loadDashboard]);

  const greetingName =
    summary?.greeting?.business_name ||
    summary?.greeting?.owner_full_name ||
    "Partner";

  const kpis = summary?.kpis;
  const recentCampaigns = summary?.recent_campaigns || [];
  const contractsSummary = summary?.contracts_summary;

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Welcome back, {greetingName}</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            Track your campaigns, discover creators, and manage collaborations.
          </p>
        </div>
        <GradientButton asChild className="h-12 rounded-xl px-6">
          <Link href="/business/campaigns/create">
            Create Campaign <ArrowRight className="size-5" />
          </Link>
        </GradientButton>
      </section>

      <ListStatus
        state={state}
        isEmpty={false}
        emptyText=""
        errorText="Your dashboard could not be loaded."
        onRetry={() => {
          setState("loading");
          void loadDashboard();
        }}
      >
        {summary?.marketplace_access && summary.marketplace_access.verification_status !== "verified" ? (
          <VerificationStatusCard
            status={summary.marketplace_access.verification_status ?? "unsubmitted"}
            action={
              <Button asChild className="h-11 rounded-xl bg-white/10 text-white hover:bg-white/15">
                <Link href="/business/profile">
                  {summary.marketplace_access.verification_status === "pending" ? "View profile" : "Go to business profile"}
                </Link>
              </Button>
            }
          />
        ) : null}

        <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <BusinessKpiCard
            icon={Target}
            label="Active Campaigns"
            value={String(kpis?.active_campaigns ?? 0)}
            trend={`${kpis?.total_campaigns ?? 0} total campaigns`}
          />
          <BusinessKpiCard
            icon={Users}
            label="Recommended Creators"
            value={String(kpis?.recommended_creators_total ?? 0)}
            trend={`${kpis?.recommended_creators_high_match ?? 0} high-match profiles`}
          />
          <BusinessKpiCard
            icon={FileText}
            label="Pending Contracts"
            value={String(kpis?.pending_contracts ?? 0)}
            trend={`${kpis?.contracts_awaiting_payment ?? 0} awaiting payment`}
          />
          <BusinessKpiCard
            icon={BarChart3}
            label="Pending Invites"
            value={String(kpis?.pending_invites ?? 0)}
            trend="Sent to creators"
          />
        </section>

        <section className="adfluence-glass rounded-2xl border border-white/10 p-5 md:p-6">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <h2 className="text-lg font-semibold">Campaign Performance</h2>
            <Button
              asChild
              variant="outline"
              className="h-10 rounded-xl border-white/15 bg-transparent px-4 text-white hover:bg-white/10"
            >
              <Link href="/business/campaigns">All campaigns</Link>
            </Button>
          </div>
          <div className="mt-4 grid gap-3 sm:grid-cols-3">
            <div className="rounded-xl bg-white/5 p-4">
              <p className="text-xs text-[#CBD5E1]">Active Campaigns</p>
              <p className="text-2xl font-semibold text-[color:var(--adfluence-orange)]">
                {kpis?.active_campaigns ?? 0}
              </p>
            </div>
            <div className="rounded-xl bg-white/5 p-4">
              <p className="text-xs text-[#CBD5E1]">Creators Needed</p>
              <p className="text-2xl font-semibold text-[color:var(--adfluence-purple-accent)]">
                {recentCampaigns.reduce((acc, c) => acc + c.creators_needed, 0)}
              </p>
            </div>
            <div className="rounded-xl bg-white/5 p-4">
              <p className="text-xs text-[#CBD5E1]">Accepted Creators</p>
              <p className="text-2xl font-semibold text-emerald-300">
                {recentCampaigns.reduce((acc, c) => acc + c.accepted_creators, 0)}
              </p>
            </div>
          </div>
        </section>

        <section className="adfluence-glass rounded-2xl border border-white/10 p-5 md:p-6">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">Recent Campaigns</h2>
            <Button asChild variant="link" className="text-sm text-[color:var(--adfluence-orange)]">
              <Link href="/business/campaigns">View All</Link>
            </Button>
          </div>

          {recentCampaigns.length === 0 ? (
            <div className="mt-6 rounded-xl border border-white/10 bg-white/5 p-6 text-center">
              <p className="text-sm text-[#CBD5E1]">No campaigns created yet.</p>
              <Button asChild className="mt-3 h-10 rounded-xl bg-white/10 text-white hover:bg-white/15">
                <Link href="/business/campaigns/create">Create Your First Campaign</Link>
              </Button>
            </div>
          ) : (
            <>
              <div className="mt-4 hidden overflow-x-auto md:block">
                <table className="w-full min-w-[680px] text-sm">
                  <thead className="text-left text-[#CBD5E1]">
                    <tr>
                      <th className="pb-3 font-medium">Campaign</th>
                      <th className="pb-3 font-medium">Niche</th>
                      <th className="pb-3 font-medium">Budget</th>
                      <th className="pb-3 font-medium">Status</th>
                      <th className="pb-3 font-medium">Creators</th>
                      <th className="pb-3 font-medium">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {recentCampaigns.map((row) => (
                      <tr key={row.id} className="border-t border-white/10">
                        <td className="py-3 font-medium">{row.title}</td>
                        <td className="py-3 text-[#CBD5E1]">{row.niche || "General"}</td>
                        <td className="py-3">{row.budget_summary || "Flexible"}</td>
                        <td className="py-3">
                          <StatusBadge status={campaignStatus(row.status)} />
                        </td>
                        <td className="py-3 text-[#CBD5E1]">
                          {row.accepted_creators} / {row.creators_needed}
                        </td>
                        <td className="py-3">
                          <Button asChild className="h-9 rounded-xl bg-white/10 text-white hover:bg-white/15">
                            {/* Only drafts and live campaigns have matches; others open the campaign. */}
                            <Link
                              href={
                                row.status === "draft" || row.status === "active"
                                  ? `/business/recommendations?campaign_id=${row.id}`
                                  : `/business/campaigns/${row.id}`
                              }
                            >
                              {row.status === "draft" || row.status === "active" ? "Matched creators" : "Open campaign"}
                            </Link>
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <div className="mt-4 space-y-3 md:hidden">
                {recentCampaigns.map((row) => (
                  <article
                    key={row.id}
                    className="rounded-xl border border-white/10 bg-white/5 p-4"
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0">
                        <p className="font-medium text-white">{row.title}</p>
                        <p className="mt-1 text-xs text-[#CBD5E1]">
                          {row.niche || "General"} · {row.budget_summary || "Flexible"}
                        </p>
                      </div>
                      <StatusBadge status={campaignStatus(row.status)} />
                    </div>
                    <div className="mt-3 flex items-center justify-between text-xs text-[#CBD5E1]">
                      <span>Creators Accepted</span>
                      <span>{row.accepted_creators} / {row.creators_needed}</span>
                    </div>
                    <Button asChild className="mt-3 h-11 w-full rounded-xl bg-white/10 text-white hover:bg-white/15">
                      {/* Only drafts and live campaigns have matches; others open the campaign. */}
                      <Link
                        href={
                          row.status === "draft" || row.status === "active"
                            ? `/business/recommendations?campaign_id=${row.id}`
                            : `/business/campaigns/${row.id}`
                        }
                      >
                        {row.status === "draft" || row.status === "active" ? "Matched creators" : "Open campaign"}
                      </Link>
                    </Button>
                  </article>
                ))}
              </div>
            </>
          )}
        </section>

        <section className="space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <h2 className="text-lg font-semibold">Recommended Influencers</h2>
            <Button
              asChild
              variant="outline"
              className="h-10 rounded-xl border-white/15 bg-transparent px-4 text-white hover:bg-white/10"
            >
              <Link href="/business/recommendations">See all recommendations</Link>
            </Button>
          </div>
          {summary?.top_recommendations && summary.top_recommendations.length > 0 ? (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
              {summary.top_recommendations.slice(0, 5).map((rec) => (
                <RecommendedInfluencerCard
                  key={rec.influencer_user_id}
                  influencer={{
                    id: rec.influencer_user_id,
                    name: rec.creator_name,
                    imageUrl: rec.profile_image_url || undefined,
                    matchScore: rec.match_score,
                    campaignId: rec.campaign_id,
                    campaignTitle: rec.campaign_title,
                  }}
                />
              ))}
            </div>
          ) : (
            <div className="rounded-2xl border border-white/10 bg-white/5 p-8 text-center">
              <Users className="mx-auto size-10 text-[#CBD5E1]/40" />
              <p className="mt-3 font-medium text-white">No creator recommendations yet</p>
              <p className="mt-1 text-sm text-[#CBD5E1]">
                Create an active marketing campaign to receive AI creator matching recommendations tailored to your niche and budget.
              </p>
              <Button
                asChild
                className="mt-4 h-10 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]"
              >
                <Link href="/business/campaigns/create">Create Campaign</Link>
              </Button>
            </div>
          )}
        </section>

        <section className="grid gap-4 lg:grid-cols-2">
          <Card className="adfluence-glass border-white/10 text-white">
            <CardHeader className="flex flex-row items-center justify-between space-y-0">
              <CardTitle>Recent Messages</CardTitle>
              <Button
                asChild
                variant="ghost"
                size="sm"
                className="h-8 text-xs text-[#CBD5E1] hover:text-white"
              >
                <Link href="/business/messages">View all</Link>
              </Button>
            </CardHeader>
            <CardContent className="space-y-3">
              {conversations.length > 0 ? (
                conversations.slice(0, 3).map((conv) => (
                  <Link
                    key={conv.id}
                    href="/business/messages"
                    className="flex items-start justify-between gap-3 rounded-xl border border-white/10 bg-white/5 p-3 transition hover:bg-white/10"
                  >
                    <div className="min-w-0">
                      <p className="font-medium text-sm text-white truncate">{conv.other_user_name}</p>
                      <p className="mt-0.5 text-xs text-[#CBD5E1] truncate">
                        {conv.last_message_preview || "Started a conversation"}
                      </p>
                    </div>
                    <div className="shrink-0 text-right">
                      <p className="text-xs text-[#CBD5E1]">
                        {conv.last_message_at
                          ? formatDate(conv.last_message_at)
                          : "New"}
                      </p>
                      {conv.unread_count > 0 ? (
                        <span className="mt-1 inline-flex items-center justify-center rounded-full bg-[color:var(--adfluence-action)] px-1.5 py-0.5 text-[10px] font-bold text-white">
                          {conv.unread_count}
                        </span>
                      ) : null}
                    </div>
                  </Link>
                ))
              ) : (
                <div className="rounded-xl border border-white/10 bg-white/5 p-6 text-center text-xs text-[#CBD5E1]">
                  <MessageSquare className="mx-auto size-6 text-[#CBD5E1]/40" />
                  <p className="mt-2 font-medium text-white">No active messages</p>
                  <p className="mt-1">Connect with creators directly through campaigns and invites.</p>
                </div>
              )}
            </CardContent>
          </Card>

          <Card className="adfluence-glass border-white/10 text-white">
            <CardHeader>
              <CardTitle>Contracts</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="rounded-xl bg-white/5 p-4">
                <p className="text-xs text-[#CBD5E1]">Pending Contracts</p>
                <p className="text-2xl font-semibold text-amber-200">
                  {contractsSummary?.pending ?? 0}
                </p>
              </div>
              <div className="rounded-xl bg-white/5 p-4">
                <p className="text-xs text-[#CBD5E1]">Active Contracts</p>
                <p className="text-2xl font-semibold text-emerald-300">
                  {contractsSummary?.active ?? 0}
                </p>
              </div>
              <div className="rounded-xl bg-white/5 p-4">
                <p className="text-xs text-[#CBD5E1]">Completed Contracts</p>
                <p className="text-2xl font-semibold text-blue-300">
                  {contractsSummary?.completed ?? 0}
                </p>
              </div>
              <Button
                asChild
                className="h-11 w-full rounded-xl bg-white/10 text-white hover:bg-white/15"
              >
                <Link href="/business/contracts">View Contracts</Link>
              </Button>
            </CardContent>
          </Card>
        </section>

        <section className="adfluence-glass rounded-2xl border border-white/10 p-5 md:p-6">
          <h2 className="text-lg font-semibold">Quick Actions</h2>
          <div className="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
            <Button
              asChild
              className="h-11 rounded-xl bg-white/10 text-white hover:bg-white/15"
            >
              <Link href="/business/campaigns/create">
                <Sparkles className="size-4" />
                Create Campaign
              </Link>
            </Button>
            <Button
              asChild
              className="h-11 rounded-xl bg-white/10 text-white hover:bg-white/15"
            >
              <Link href="/business/recommendations">
                <Users className="size-4" />
                Browse Influencers
              </Link>
            </Button>
            <Button
              asChild
              className="h-11 rounded-xl bg-white/10 text-white hover:bg-white/15"
            >
              <Link href="/business/messages">
                <MessageSquare className="size-4" />
                Messages
              </Link>
            </Button>
            <Button
              asChild
              className="h-11 rounded-xl bg-white/10 text-white hover:bg-white/15"
            >
              <Link href="/business/contracts">
                <FileText className="size-4" />
                Contracts
              </Link>
            </Button>
          </div>
        </section>
      </ListStatus>
    </div>
  );
}
