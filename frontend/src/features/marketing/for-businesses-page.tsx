import { PageHeading } from "@/components/common/page-heading";
import { StatCard } from "@/components/dashboard/stat-card";
import { businessStats } from "@/data/stats";

export function ForBusinessesPage() {
  return (
    <section className="mx-auto max-w-6xl space-y-6 px-4 py-10">
      <PageHeading
        title="For Businesses"
        description="Discover creators, run campaigns, and manage collaborations."
      />
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {businessStats.map((item) => (
          <StatCard key={item.id} item={item} />
        ))}
      </div>
    </section>
  );
}
