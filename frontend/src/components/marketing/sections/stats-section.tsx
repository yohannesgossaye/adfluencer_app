import { Container } from "@/components/common/container";

export function StatsSection() {
  const stats = [
    { value: "500+", label: "Active Influencers" },
    { value: "200+", label: "Brands Trust Us" },
    { value: "1,000+", label: "Campaigns Completed" },
    { value: "98%", label: "Client Satisfaction" },
  ];

  return (
    <section className="bg-[color:var(--adfluence-surface)]/50 py-20 text-white">
      <Container>
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
        {stats.map((item) => (
          <div key={item.label} className="text-center">
            <div className="mb-2 text-5xl font-bold text-[color:var(--adfluence-orange)]">
              {item.value}
            </div>
            <div className="text-[#CBD5E1]">{item.label}</div>
          </div>
        ))}
        </div>
      </Container>
    </section>
  );
}
