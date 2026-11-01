import type { ReactNode } from "react";

import { ChatProvider } from "@/components/chat/chat-provider";
import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { businessNav } from "@/lib/navigation";

export default function BusinessLayout({ children }: { children: ReactNode }) {
  return (
    <ChatProvider>
      <DashboardShell
        sidebarTitle="Business workspace"
        navItems={businessNav}
        headerTitle="Business Portal"
        headerSubtitle="Manage campaigns and creator collaborations"
      >
        {children}
      </DashboardShell>
    </ChatProvider>
  );
}
