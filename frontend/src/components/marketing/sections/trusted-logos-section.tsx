import { Container } from "@/components/common/container";

export function TrustedLogosSection() {
  const logos = ["Dashen", "Abyssinia", "Ride", "Hahu", "Yaya", "Kifiya"];

  return (
    <section className="border-t border-white/10 bg-[color:var(--adfluence-navy)] py-14">
      <Container>
        <p className="text-center text-sm text-slate-300">
          TRUSTED BY LEADING ETHIOPIAN BRANDS
        </p>
        <div className="mt-8 flex flex-wrap items-center justify-center gap-12 opacity-40">
          {logos.map((logo) => (
            <div key={logo} className="text-2xl font-bold text-white">
              {logo}
            </div>
          ))}
        </div>
      </Container>
    </section>
  );
}
