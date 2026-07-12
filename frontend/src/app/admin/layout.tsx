"use client";

import type { ReactNode } from "react";
import { usePathname } from "next/navigation";

import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { adminNav } from "@/lib/navigation";

export default function AdminLayout({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  if (pathname === "/admin/login") {
    return (
      <div className="relative flex min-h-svh flex-col items-center justify-center overflow-hidden bg-[color:var(--adfluence-bg)] p-4 sm:p-6">
        <div className="absolute -left-24 top-8 h-72 w-72 rounded-full bg-[#0B2240]/40 blur-[110px]" />
        <div className="absolute right-[-12%] top-[15%] h-80 w-80 rounded-full bg-[color:var(--adfluence-orange)]/8 blur-[120px]" />
        <div className="relative z-10 w-full max-w-2xl">{children}</div>
      </div>
    );
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
