import Link from "next/link";

import { Button } from "@/components/ui/button";

const businessCampaigns = [
  {
    id: "bcmp-1",
    title: "Betu Cosmetics Launch",
    brand: "Betu Cosmetics",
    platform: "Instagram",
    budget: "ETB 280,000",
    status: "Active",
    progress: 72,
  },
  {
    id: "bcmp-2",
    title: "Coffee House Awareness",
    brand: "Habesha Coffee House",
    platform: "TikTok",
    budget: "ETB 190,000",
    status: "Pending",
    progress: 45,
  },
  {
    id: "bcmp-3",
    title: "Ride App Creator Push",
    brand: "Ride Ethiopia",
    platform: "YouTube",
    budget: "ETB 340,000",
    status: "Completed",
    progress: 100,
  },
  {
    id: "bcmp-4",
    title: "Abyssinia Restaurant Promo",
    brand: "Abyssinia Restaurant",
    platform: "Instagram",
    budget: "ETB 150,000",
    status: "Active",
    progress: 64,
  },
] as const;

export function CampaignsList() {
  return (
    <div className="space-y-6 rounded-2xl bg-[#0B1A2B] p-4 text-white md:p-6">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-white">Campaigns</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            Manage campaign lifecycle from draft to completion.
          </p>
        </div>
        <Button
          asChild
          className="h-11 rounded-xl bg-gradient-to-r from-[#FF7A00] to-[#7C3AED] px-5 text-white hover:opacity-95"
        >
          <Link href="/business/campaigns/create">Create Campaign</Link>
        </Button>
      </section>

      <div className="grid gap-4 md:grid-cols-2">
        {businessCampaigns.map((campaign) => (
          <article
            key={campaign.id}
            className="adfluence-glass rounded-2xl border border-white/10 p-5"
          >
            <div className="flex items-start justify-between gap-3">
              <div>
                <h2 className="text-lg font-semibold text-white">{campaign.title}</h2>
                <p className="mt-1 text-sm text-[#CBD5E1]">{campaign.brand}</p>
              </div>
              <span
                className={`rounded-full px-3 py-1 text-xs ${
                  campaign.status === "Active"
                    ? "bg-emerald-500/15 text-emerald-300"
                    : campaign.status === "Pending"
                      ? "bg-amber-500/15 text-amber-200"
                      : "bg-blue-500/15 text-blue-300"
                }`}
              >
                {campaign.status}
              </span>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-3 text-sm text-[#CBD5E1]">
              <p>Platform: {campaign.platform}</p>
              <p>Budget: {campaign.budget}</p>
            </div>

            <div className="mt-4">
              <p className="mb-2 text-xs text-[#CBD5E1]">Progress</p>
              <div className="h-2 rounded-full bg-white/10">
                <div
                  className="h-2 rounded-full bg-gradient-to-r from-[#FF7A00] to-[#7C3AED]"
                  style={{ width: `${campaign.progress}%` }}
                />
              </div>
            </div>

            <Button className="mt-4 h-10 rounded-xl bg-white/10 text-white hover:bg-white/15">
              View
            </Button>
          </article>
        ))}
      </div>
    </div>
  );
}
