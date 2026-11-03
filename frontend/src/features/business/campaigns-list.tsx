"use client";

import Link from "next/link";
import { useCallback } from "react";

import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import { usePagedList } from "@/hooks/use-paged-list";
import { businessApi, type CampaignListItem } from "@/lib/api/business";

function getStatusBadge(status: string) {
  const s = status.toLowerCase();
  if (s === "active") return "bg-emerald-500/15 text-emerald-300";
  if (s === "draft") return "bg-amber-500/15 text-amber-200";
  if (s === "paused") return "bg-rose-500/15 text-rose-300";
  if (s === "completed") return "bg-blue-500/15 text-blue-300";
  return "bg-white/10 text-[#CBD5E1]";
}

const STATUS_LABEL: Record<string, string> = {
  draft: "Draft",
  active: "Live",
  paused: "Paused by Adfluence",
  completed: "Completed",
  cancelled: "Cancelled",
};

export function CampaignsList() {
  const fetchPage = useCallback(async (page: number) => (await businessApi.listCampaigns({ page })).data, []);
  const list = usePagedList(fetchPage);
  const campaigns: CampaignListItem[] = list.items;

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-white">Campaigns</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            Manage campaign lifecycle from draft to completion.
          </p>
        </div>
        <Button
          asChild
          className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]"
        >
          <Link href="/business/campaigns/create">Create Campaign</Link>
        </Button>
      </section>

      <ListStatus
        state={list.state}
        isEmpty={campaigns.length === 0}
        errorText="Your campaigns could not be loaded."
        onRetry={list.retry}
        emptyText={
          <>
            <p className="text-base text-white">No campaigns yet.</p>
            <p className="mt-1 text-sm text-[#CBD5E1]">Create a new campaign to start finding matched creators.</p>
            <Button
              asChild
              className="mt-4 h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]"
            >
              <Link href="/business/campaigns/create">Create Campaign</Link>
            </Button>
          </>
        }
      >
          <div className="grid gap-4 md:grid-cols-2">
            {campaigns.map((campaign) => (
              <article
                key={campaign.id}
                className="adfluence-glass rounded-2xl border border-white/10 p-5"
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <h2 className="text-lg font-semibold text-white">
                      <Link href={`/business/campaigns/${campaign.id}`} className="hover:underline">
                        {campaign.title}
                      </Link>
                    </h2>
                    <p className="mt-1 text-sm text-[#CBD5E1]">{campaign.business_name}</p>
                  </div>
                  <span
                    className={`shrink-0 rounded-full px-3 py-1 text-xs ${getStatusBadge(campaign.status)}`}
                  >
                    {STATUS_LABEL[campaign.status] ?? campaign.status}
                  </span>
                </div>

                <div className="mt-4 grid grid-cols-2 gap-3 text-sm text-[#CBD5E1]">
                  <p>Platform: <span className="text-white capitalize">{campaign.primary_platform || "Multi"}</span></p>
                  <p>Budget: <span className="text-white">{campaign.budget_summary || "Flexible"}</span></p>
                </div>

                <div className="mt-4">
                  <div className="mb-2 flex items-center justify-between text-xs text-[#CBD5E1]">
                    <span>Creators Accepted</span>
                    <span className="text-white">{campaign.invites_accepted} / {campaign.creators_needed}</span>
                  </div>
                  <div className="h-2 rounded-full bg-white/10">
                    <div
                      className="h-2 rounded-full bg-[color:var(--adfluence-action)]"
                      style={{
                        width: `${Math.min(
                          100,
                          campaign.creators_needed > 0
                            ? Math.round((campaign.invites_accepted / campaign.creators_needed) * 100)
                            : 0
                        )}%`,
                      }}
                    />
                  </div>
                </div>

                <div className="mt-5 flex flex-wrap items-center justify-end gap-2">
                  {campaign.status === "active" ? (
                    <Button asChild className="h-10 rounded-xl bg-white/10 text-white hover:bg-white/15">
                      <Link href={`/business/recommendations?campaign_id=${campaign.id}`}>Find creators</Link>
                    </Button>
                  ) : null}
                  <Button
                    asChild
                    className="h-10 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]"
                  >
                    <Link href={`/business/campaigns/${campaign.id}`}>
                      {campaign.status === "draft" ? "Review & publish" : "Open campaign"}
                    </Link>
                  </Button>
                </div>
              </article>
            ))}
          </div>
        <LoadMore
          shown={campaigns.length}
          total={list.total}
          hasMore={list.hasMore}
          loading={list.loadingMore}
          onLoadMore={list.loadMore}
        />
      </ListStatus>
    </div>
  );
}

