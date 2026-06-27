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

const performanceBars = [35, 54, 41, 70, 58, 79, 67, 84, 72, 91];

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

export function BusinessDashboard() {
  return (
    <div className="space-y-6 rounded-2xl bg-[#0B1A2B] p-4 text-white md:p-6">
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
        <div className="mt-5 flex h-44 items-end gap-2 rounded-xl bg-white/5 p-4">
          {performanceBars.map((bar, i) => (
            <div
              key={i}
              className="flex-1 rounded-t bg-gradient-to-t from-[#FF7A00] to-[#7C3AED]"
              style={{ height: `${bar}%` }}
            />
          ))}
        </div>
        <div className="mt-4 grid gap-3 sm:grid-cols-3">
          <div className="rounded-xl bg-white/5 p-4">
            <p className="text-xs text-[#CBD5E1]">Reach</p>
            <p className="text-2xl font-semibold text-[#FF7A00]">245K</p>
          </div>
          <div className="rounded-xl bg-white/5 p-4">
            <p className="text-xs text-[#CBD5E1]">Engagement</p>
            <p className="text-2xl font-semibold text-[#7C3AED]">38.4K</p>
          </div>
          <div className="rounded-xl bg-white/5 p-4">
            <p className="text-xs text-[#CBD5E1]">Conversion Lift</p>
            <p className="text-2xl font-semibold text-emerald-300">+24%</p>
          </div>
        </div>
      </section>

      <section className="adfluence-glass rounded-2xl border border-white/10 p-5 md:p-6">
        <h2 className="text-lg font-semibold">Active Campaigns</h2>
        <div className="mt-4 overflow-x-auto">
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
                      className={`rounded-full px-3 py-1 text-xs ${
                        row.status === "Active"
                          ? "bg-emerald-500/15 text-emerald-300"
                          : row.status === "Pending"
                            ? "bg-amber-500/15 text-amber-200"
                            : "bg-blue-500/15 text-blue-300"
                      }`}
                    >
                      {row.status}
                    </span>
                  </td>
                  <td className="py-3">
                    <div className="w-32 rounded-full bg-white/10">
                      <div
                        className="h-2 rounded-full bg-gradient-to-r from-[#FF7A00] to-[#7C3AED]"
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
                    <span className="mt-1 inline-block h-2 w-2 rounded-full bg-[#FF7A00]" />
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
