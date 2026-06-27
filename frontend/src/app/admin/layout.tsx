"use client";

import type { ReactNode } from "react";
import { usePathname } from "next/navigation";

import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { adminNav } from "@/lib/navigation";

export default function AdminLayout({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  if (pathname === "/admin/login") {
    return <div className="bg-muted/30 min-h-svh p-6">{children}</div>;
  }

  return (
    <DashboardShell
      sidebarTitle="Admin workspace"
      navItems={adminNav}
      headerTitle="Admin Portal"
      headerSubtitle="Moderate platform users and campaigns"
    >
      {children}
    </DashboardShell>
  );
}
