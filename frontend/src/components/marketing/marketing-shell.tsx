import type { ReactNode } from "react";

import { MarketingNavbar } from "@/components/marketing/marketing-navbar";

export function MarketingShell({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-[color:var(--adfluence-navy)] text-white">
      <MarketingNavbar />
      <main>{children}</main>
    </div>
  );
}
