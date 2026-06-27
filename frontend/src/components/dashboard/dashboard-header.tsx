import type { ReactNode } from "react";

export function DashboardHeader({
  title,
  subtitle,
  leading,
}: {
  title: string;
  subtitle: string;
  leading?: ReactNode;
}) {
  return (
    <header className="border-b border-white/10 bg-[rgba(11,26,43,0.95)] px-4 py-4 backdrop-blur-xl md:px-6">
      <div className="flex items-start gap-3">
        {leading ? <div className="mt-0.5 md:hidden">{leading}</div> : null}
        <div>
          <h1 className="text-xl font-semibold text-white">{title}</h1>
          <p className="text-sm text-[#CBD5E1]">{subtitle}</p>
        </div>
      </div>
    </header>
  );
}
