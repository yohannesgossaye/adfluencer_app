import type { LucideIcon } from "lucide-react";

type BusinessKpiCardProps = {
  icon: LucideIcon;
  label: string;
  value: string;
  trend: string;
};

export function BusinessKpiCard({
  icon: Icon,
  label,
  value,
  trend,
}: BusinessKpiCardProps) {
  return (
    <article className="adfluence-glass rounded-2xl border border-white/10 p-5 text-white">
      <div className="mb-4 flex items-center justify-between">
        <span className="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED] text-white">
          <Icon className="size-5" />
        </span>
      </div>
      <p className="text-3xl font-bold">{value}</p>
      <p className="mt-1 text-sm text-[#CBD5E1]">{label}</p>
      <p className="mt-3 text-xs text-[#CBD5E1]">{trend}</p>
    </article>
  );
}

