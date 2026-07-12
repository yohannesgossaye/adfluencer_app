import Link from "next/link";
import {
  ArrowRight,
  BarChart3,
  FileText,
  MessageSquare,
  Sparkles,
  Target,
  Users,
} from "lucide-react";

import { RecommendedInfluencerCard } from "@/components/cards/recommended-influencer-card";
import { GradientButton } from "@/components/common/gradient-button";
import { BusinessKpiCard } from "@/components/dashboard/business-kpi-card";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { contracts } from "@/data/contracts";
import { influencers } from "@/data/influencers";
import { messages } from "@/data/messages";

const activeCampaignRows = [
  {
    name: "Betu Cosmetics Launch",
    niche: "Beauty",
    budget: "ETB 280K",
    status: "Active",
    progress: 72,
  },
  {
    name: "Coffee House Awareness",
    niche: "Food",
    budget: "ETB 190K",
    status: "Pending",
    progress: 45,
  },
  {
    name: "Ride App Creator Push",
    niche: "Lifestyle",
    budget: "ETB 340K",
    status: "Completed",
    progress: 100,
  },
];

function campaignStatusClass(status: string) {
  if (status === "Active") return "bg-emerald-500/15 text-emerald-300";
  if (status === "Pending") return "bg-amber-500/15 text-amber-200";
  return "bg-blue-500/15 text-blue-300";
}

export function BusinessDashboard() {
  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Welcome back, Yohannes</h1>
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

      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <BusinessKpiCard
          icon={Target}
          label="Active Campaigns"
          value="3"
          trend="+1 this week"
        />
        <BusinessKpiCard
          icon={Users}
          label="Recommended Creators"
          value="24"
          trend="12 high-match profiles"
        />
        <BusinessKpiCard
          icon={FileText}
          label="Pending Contracts"
          value="5"
          trend="2 require approval"
        />
        <BusinessKpiCard
          icon={BarChart3}
          label="Total Reach"
          value="245K"
          trend="+18% vs last month"
        />
      </section>

      <section className="adfluence-glass rounded-2xl border border-white/10 p-5 md:p-6">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <h2 className="text-lg font-semibold">Campaign Performance</h2>
          <Button
            variant="outline"
            className="h-10 rounded-xl border-white/15 bg-transparent px-4 text-white hover:bg-white/10"
          >
            Last 30 days
          </Button>
        </div>
        <div className="mt-4 grid gap-3 sm:grid-cols-3">
          <div className="rounded-xl bg-white/5 p-4">
            <p className="text-xs text-[#CBD5E1]">Reach</p>
            <p className="text-2xl font-semibold text-[color:var(--adfluence-orange)]">245K</p>
          </div>
          <div className="rounded-xl bg-white/5 p-4">
            <p className="text-xs text-[#CBD5E1]">Engagement</p>
            <p className="text-2xl font-semibold text-[color:var(--adfluence-purple-accent)]">38.4K</p>
          </div>
          <div className="rounded-xl bg-white/5 p-4">
            <p className="text-xs text-[#CBD5E1]">Conversion Lift</p>
            <p className="text-2xl font-semibold text-emerald-300">+24%</p>
          </div>
        </div>
      </section>

      <section className="adfluence-glass rounded-2xl border border-white/10 p-5 md:p-6">
        <h2 className="text-lg font-semibold">Active Campaigns</h2>
        <div className="mt-4 hidden overflow-x-auto md:block">
          <table className="w-full min-w-[680px] text-sm">
            <thead className="text-left text-[#CBD5E1]">
              <tr>
                <th className="pb-3 font-medium">Campaign</th>
                <th className="pb-3 font-medium">Niche</th>
                <th className="pb-3 font-medium">Budget</th>
                <th className="pb-3 font-medium">Status</th>
                <th className="pb-3 font-medium">Progress</th>
                <th className="pb-3 font-medium">Action</th>
              </tr>
            </thead>
            <tbody>
              {activeCampaignRows.map((row) => (
                <tr key={row.name} className="border-t border-white/10">
                  <td className="py-3">{row.name}</td>
                  <td className="py-3">{row.niche}</td>
                  <td className="py-3">{row.budget}</td>
                  <td className="py-3">
                    <span
                      className={`rounded-full px-3 py-1 text-xs ${campaignStatusClass(row.status)}`}
                    >
                      {row.status}
                    </span>
                  </td>
                  <td className="py-3">
                    <div className="w-32 rounded-full bg-white/10">
                      <div
                        className="h-2 rounded-full bg-[color:var(--adfluence-orange)]"
                        style={{ width: `${row.progress}%` }}
                      />
                    </div>
                  </td>
                  <td className="py-3">
                    <Button className="h-9 rounded-xl bg-white/10 text-white hover:bg-white/15">
                      View
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="mt-4 space-y-3 md:hidden">
          {activeCampaignRows.map((row) => (
            <article
              key={row.name}
              className="rounded-xl border border-white/10 bg-white/5 p-4"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="font-medium text-white">{row.name}</p>
                  <p className="mt-1 text-xs text-[#CBD5E1]">
                    {row.niche} · {row.budget}
                  </p>
                </div>
                <span
                  className={`shrink-0 rounded-full px-3 py-1 text-xs ${campaignStatusClass(row.status)}`}
                >
                  {row.status}
                </span>
              </div>
              <div className="mt-3">
                <div className="mb-1 flex items-center justify-between text-xs text-[#CBD5E1]">
                  <span>Progress</span>
                  <span>{row.progress}%</span>
                </div>
                <div className="rounded-full bg-white/10">
                  <div
                    className="h-2 rounded-full bg-[color:var(--adfluence-orange)]"
                    style={{ width: `${row.progress}%` }}
                  />
                </div>
              </div>
              <Button className="mt-3 h-11 w-full rounded-xl bg-white/10 text-white hover:bg-white/15 md:h-9">
                View
              </Button>
            </article>
          ))}
        </div>
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
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
          {influencers.slice(0, 5).map((influencer) => (
            <RecommendedInfluencerCard key={influencer.id} influencer={influencer} />
          ))}
        </div>
      </section>

      <section className="grid gap-4 lg:grid-cols-2">
        <Card className="adfluence-glass border-white/10 text-white">
          <CardHeader>
            <CardTitle>Recent Messages</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {messages.slice(0, 3).map((message) => (
              <div
                key={message.id}
                className="flex items-start justify-between gap-3 rounded-xl border border-white/10 bg-white/5 p-3"
              >
                <div>
                  <p className="text-sm font-medium">{message.from}</p>
                  <p className="text-xs text-[#CBD5E1]">{message.preview}</p>
                </div>
                <div className="text-right">
                  <p className="text-xs text-[#CBD5E1]">
                    {new Date(message.timestamp).toLocaleDateString()}
                  </p>
                  {message.unread ? (
                    <span className="mt-1 inline-block h-2 w-2 rounded-full bg-[color:var(--adfluence-orange)]" />
                  ) : null}
                </div>
              </div>
            ))}
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
                {contracts.filter((c) => c.status === "pending").length}
              </p>
            </div>
            <div className="rounded-xl bg-white/5 p-4">
              <p className="text-xs text-[#CBD5E1]">Active Contracts</p>
              <p className="text-2xl font-semibold text-emerald-300">
                {contracts.filter((c) => c.status === "active").length}
              </p>
            </div>
            <div className="rounded-xl bg-white/5 p-4">
              <p className="text-xs text-[#CBD5E1]">Completed Contracts</p>
              <p className="text-2xl font-semibold text-blue-300">
                {contracts.filter((c) => c.status === "completed").length}
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
    </div>
  );
}
