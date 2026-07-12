import { Container } from "@/components/common/container";

export function StatsSection() {
  const stats = [
    { value: "500+", label: "Active Influencers" },
    { value: "200+", label: "Brands Trust Us" },
    { value: "1,000+", label: "Campaigns Completed" },
    { value: "98%", label: "Client Satisfaction" },
  ];

  return (
    <section className="bg-[color:var(--adfluence-surface)]/50 py-12 text-white sm:py-20">
      <Container>
        <div className="grid grid-cols-2 gap-4 sm:gap-8 md:grid-cols-4">
        {stats.map((item) => (
          <div key={item.label} className="text-center">
            <div className="mb-1 text-3xl font-bold text-[color:var(--adfluence-orange)] sm:mb-2 sm:text-5xl">
              {item.value}
            </div>
            <div className="text-xs text-[#CBD5E1] sm:text-base">{item.label}</div>
          </div>
        ))}
        </div>
      </Container>
    </section>
  );
}
