import { Container } from "@/components/common/container";

const stats = [
  { value: "500+", label: "Total Creators" },
  { value: "ETB 12M+", label: "Campaign Payouts" },
  { value: "200+", label: "Brands Hiring" },
  { value: "3.4x", label: "Avg Creator Growth" },
];

export function InfluencerTrustStatsSection() {
  return (
    <section className="bg-[color:var(--adfluence-navy-dark)]/50 py-12 text-white sm:py-20">
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

