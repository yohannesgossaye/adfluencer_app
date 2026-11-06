"use client";

import * as React from "react";
import { Toaster } from "sonner";

import { AuthProvider } from "@/components/providers/auth-provider";
import { SessionGuard } from "@/components/auth/session-guard";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryProvider } from "@/lib/query-client";

export function AppProviders({ children }: { children: React.ReactNode }) {
  return (
    <QueryProvider>
      <AuthProvider>
        <TooltipProvider delayDuration={0}>
          {children}
          <SessionGuard />
          <Toaster theme="dark" richColors position="top-center" />
        </TooltipProvider>
      </AuthProvider>
    </QueryProvider>
  );
}
