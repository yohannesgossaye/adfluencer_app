"use client";

import * as React from "react";
import { Toaster } from "sonner";

import { ActiveThemeProvider } from "@/components/providers/active-theme-provider";
import { AuthProvider } from "@/components/providers/auth-provider";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryProvider } from "@/lib/query-client";

type AppProvidersProps = {
  children: React.ReactNode;
  defaultTheme?: "light" | "dark" | "system";
  themeScale?: string;
};

export function AppProviders({
  children,
  defaultTheme = "system",
  themeScale = "100",
}: AppProvidersProps) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme={defaultTheme}
      enableSystem
      disableTransitionOnChange
    >
      <QueryProvider>
        <AuthProvider>
          <ActiveThemeProvider themeScale={themeScale}>
            <TooltipProvider delayDuration={0}>
              {children}
              <Toaster richColors position="top-center" />
            </TooltipProvider>
          </ActiveThemeProvider>
        </AuthProvider>
      </QueryProvider>
    </ThemeProvider>
  );
}
