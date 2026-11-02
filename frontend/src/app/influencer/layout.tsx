import type { ReactNode } from "react";

import { RoleGuard } from "@/components/auth/role-guard";
import { ChatProvider } from "@/components/chat/chat-provider";
import { InfluencerAccountProvider } from "@/components/influencer/account-status";
import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { influencerNav } from "@/lib/navigation";

export default function InfluencerLayout({ children }: { children: ReactNode }) {
  return (
    <RoleGuard role="influencer">
      <InfluencerAccountProvider>
        <ChatProvider>
          <DashboardShell
            sidebarTitle="Influencer workspace"
            navItems={influencerNav}
            headerTitle="Influencer Portal"
            headerSubtitle="Manage profile, invites, and collaborations"
          >
            {children}
          </DashboardShell>
        </ChatProvider>
      </InfluencerAccountProvider>
    </RoleGuard>
  );
}
