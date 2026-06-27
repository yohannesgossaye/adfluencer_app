import { PageHeading } from "@/components/common/page-heading";
import { StatCard } from "@/components/dashboard/stat-card";
import { influencerStats } from "@/data/stats";

export function ForInfluencersPage() {
  return (
    <section className="mx-auto max-w-6xl space-y-6 px-4 py-10">
      <PageHeading
        title="For Influencers"
        description="Showcase your profile and receive invitations from brands."
      />
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {influencerStats.map((item) => (
          <StatCard key={item.id} item={item} />
        ))}
      </div>
    </section>
  );
}
