import { Quote } from "lucide-react";

import { Container } from "@/components/common/container";

const logos = ["Dashen", "Abyssinia", "Ride", "Hahu", "Yaya", "Kifiya"];

const stats = [
  { value: "500+", label: "Creators" },
  { value: "200+", label: "Campaigns" },
  { value: "98%", label: "Satisfaction" },
  { value: "3.8x", label: "Avg. ROI" },
];

export function TrustResultsSection() {
  return (
    <section className="py-20 bg-[color:var(--adfluence-navy-dark)]/50">
      <Container>
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-white sm:text-4xl">
            Trusted. Measurable. Repeatable.
          </h2>
          <p className="mt-3 text-[#CBD5E1]">
            Built for teams that need real outcomes — not guesswork.
          </p>
        </div>

        <div className="flex flex-wrap items-center justify-center gap-10 opacity-40">
          {logos.map((logo) => (
            <div key={logo} className="text-2xl font-bold text-white">
              {logo}
            </div>
          ))}
        </div>

        <div className="mt-10 grid gap-6 md:grid-cols-2">
          <div className="adfluence-glass rounded-2xl border border-white/15 p-8">
            <div className="flex items-center gap-3">
              <span className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-[color:var(--adfluence-purple)] text-white">
                <Quote className="size-5" />
              </span>
              <p className="text-sm text-[#CBD5E1]">Testimonial</p>
            </div>
            <p className="mt-5 text-lg leading-8 text-white">
              “We moved from scattered DMs to a clear workflow. The creator
              matches were accurate, and the reporting made ROI visible.”
            </p>
            <p className="mt-5 text-sm text-[#CBD5E1]">
              Marketing Lead • Addis Ababa
            </p>
          </div>

          <div className="grid grid-cols-2 gap-4">
            {stats.map((s) => (
              <div
                key={s.label}
                className="adfluence-glass rounded-2xl border border-white/15 p-6 text-center"
              >
                <p className="text-3xl font-bold text-[color:var(--adfluence-orange)] sm:text-4xl">
                  {s.value}
                </p>
                <p className="mt-1 text-xs text-[#CBD5E1] sm:mt-2 sm:text-sm">{s.label}</p>
              </div>
            ))}
          </div>
        </div>
      </Container>
    </section>
  );
}

