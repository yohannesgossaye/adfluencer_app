import { BadgeCheck, TrendingUp } from "lucide-react";

import { Container } from "@/components/common/container";

function MetricCard({
  label,
  value,
  accent,
}: {
  label: string;
  value: string;
  accent: "orange" | "purple" | "green" | "amber";
}) {
  const valueClass =
    accent === "orange"
      ? "text-[#FF7A00]"
      : accent === "purple"
        ? "text-[#7C3AED]"
        : accent === "green"
          ? "text-emerald-300"
          : "text-amber-200";

  return (
    <div className="rounded-xl bg-white/5 p-4">
      <div className={`text-2xl font-bold ${valueClass} mb-1`}>{value}</div>
      <div className="text-xs text-[#CBD5E1]">{label}</div>
    </div>
  );
}

export function BusinessDashboardPreviewSection() {
  return (
    <section className="py-20">
      <Container>
        <div className="grid gap-10 lg:grid-cols-2 lg:items-center">
          <div>
            <h2 className="text-3xl font-bold text-white sm:text-4xl">
              A dashboard built for ROI
            </h2>
            <p className="mt-4 text-[#CBD5E1]">
              Track your pipeline, creator collaborations, and campaign results in
              one place — with clear status and performance signals.
            </p>
            <div className="mt-6 grid gap-3 sm:grid-cols-2">
              {[
                "Campaign status and approvals",
                "Creator shortlist and invites",
                "Engagement and conversion metrics",
                "ROI and spend visibility",
              ].map((t) => (
                <div key={t} className="flex items-center gap-2 text-sm text-[#CBD5E1]">
                  <BadgeCheck className="size-4 text-emerald-300" />
                  {t}
                </div>
              ))}
            </div>
          </div>

          <div className="adfluence-glass rounded-2xl border border-white/20 p-6">
            <div className="flex items-center justify-between pb-4 border-b border-white/10">
              <h4 className="font-semibold text-white">Campaign Dashboard</h4>
              <span className="inline-flex items-center rounded-full bg-emerald-500/15 px-3 py-1 text-xs font-medium text-emerald-300 border border-emerald-500/30">
                Live
              </span>
            </div>

            <div className="grid grid-cols-2 gap-4 mt-4">
              <MetricCard label="Total Impressions" value="2.4M" accent="orange" />
              <MetricCard label="Engagement Rate" value="8.2%" accent="purple" />
              <MetricCard label="Conversions" value="12K" accent="green" />
              <MetricCard label="ROI" value="3.8x" accent="amber" />
            </div>

            <div className="mt-4 h-32 bg-white/5 rounded-xl flex items-end p-4 gap-2">
              {[40, 65, 45, 80, 60, 90, 75].map((h, i) => (
                <div
                  key={i}
                  className="flex-1 bg-gradient-to-t from-[#FF7A00] to-[#7C3AED] rounded-t"
                  style={{ height: `${h}%` }}
                />
              ))}
            </div>

            <div className="mt-4 grid gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <p className="text-xs text-[#CBD5E1]">Top creators</p>
                <div className="mt-3 space-y-2">
                  {[
                    { name: "Addis Food Trails", tag: "Food • TikTok" },
                    { name: "Selam Beauty Studio", tag: "Beauty • Instagram" },
                    { name: "Miki Fit Ethiopia", tag: "Fitness • YouTube" },
                  ].map((c) => (
                    <div
                      key={c.name}
                      className="flex items-center justify-between rounded-xl bg-white/5 px-3 py-2"
                    >
                      <div>
                        <p className="text-sm font-medium text-white">{c.name}</p>
                        <p className="text-xs text-[#CBD5E1]">{c.tag}</p>
                      </div>
                      <span className="text-xs text-white/70">View</span>
                    </div>
                  ))}
                </div>
              </div>
              <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <p className="text-xs text-[#CBD5E1]">Campaign status</p>
                <div className="mt-3 space-y-2">
                  {[
                    { label: "Draft", value: "2" },
                    { label: "Active", value: "8" },
                    { label: "Completed", value: "14" },
                  ].map((s) => (
                    <div
                      key={s.label}
                      className="flex items-center justify-between rounded-xl bg-white/5 px-3 py-2"
                    >
                      <p className="text-sm text-white">{s.label}</p>
                      <p className="text-sm font-semibold text-white">{s.value}</p>
                    </div>
                  ))}
                </div>
                <div className="mt-3 flex items-center gap-2 text-xs text-[#CBD5E1]">
                  <TrendingUp className="size-4 text-emerald-300" />
                  Performance improving week over week
                </div>
              </div>
            </div>
          </div>
        </div>
      </Container>
    </section>
  );
}

