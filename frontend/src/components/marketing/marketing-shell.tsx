"use client";

import type { ReactNode } from "react";
import { Suspense } from "react";

import { MarketingAuthModal } from "@/components/marketing/marketing-auth-modal";
import { MarketingAuthProvider } from "@/components/marketing/marketing-auth-provider";
import { MarketingNavbar } from "@/components/marketing/marketing-navbar";

function MarketingShellContent({ children }: { children: ReactNode }) {
  return (
    <MarketingAuthProvider>
      <div className="min-h-screen bg-[color:var(--adfluence-bg)] text-white">
        <MarketingNavbar />
        <main>{children}</main>
        <MarketingAuthModal />
      </div>
    </MarketingAuthProvider>
  );
}

export function MarketingShell({ children }: { children: ReactNode }) {
  return (
    <Suspense fallback={<div className="min-h-screen bg-[color:var(--adfluence-bg)]" />}>
      <MarketingShellContent>{children}</MarketingShellContent>
    </Suspense>
  );
}
