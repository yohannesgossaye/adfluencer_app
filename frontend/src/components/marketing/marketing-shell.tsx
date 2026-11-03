"use client";

import { Suspense, type ReactNode } from "react";

import { MarketingAuthModal } from "@/components/marketing/marketing-auth-modal";
import { MarketingAuthProvider } from "@/components/marketing/marketing-auth-provider";
import { MarketingFooter } from "@/components/marketing/marketing-footer";
import { MarketingNavbar } from "@/components/marketing/marketing-navbar";

export function MarketingShell({ children }: { children: ReactNode }) {
  return (
    <MarketingAuthProvider>
      <div className="min-h-screen bg-[color:var(--adfluence-bg)] text-white">
        <MarketingNavbar />
        <main>{children}</main>
        <MarketingFooter />
        {/* Only the dialog reads the URL, so pages are rendered on the server (good for Google). */}
        <Suspense fallback={null}>
          <MarketingAuthModal />
        </Suspense>
      </div>
    </MarketingAuthProvider>
  );
}
