import type { ReactNode } from "react";

import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { influencerNav } from "@/lib/navigation";

export default function InfluencerLayout({ children }: { children: ReactNode }) {
  return (
    <DashboardShell
      sidebarTitle="Influencer workspace"
      navItems={influencerNav}
      headerTitle="Influencer Portal"
      headerSubtitle="Manage profile, invites, and collaborations"
    >
      {children}
    </DashboardShell>
  );
}
