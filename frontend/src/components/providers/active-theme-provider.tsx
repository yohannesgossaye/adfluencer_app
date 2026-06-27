"use client";

import { useTheme } from "next-themes";
import * as React from "react";

type ActiveThemeProviderProps = {
  children: React.ReactNode;
  /** Optional scale token from server (e.g. cookie): `100` | `90` | `95` */
  themeScale?: string;
};

export function ActiveThemeProvider({
  children,
  themeScale = "100",
}: ActiveThemeProviderProps) {
  const { resolvedTheme } = useTheme();

  React.useEffect(() => {
    document.documentElement.dataset.theme = resolvedTheme ?? "";
    document.documentElement.dataset.themeScale = themeScale;
  }, [resolvedTheme, themeScale]);

  return <>{children}</>;
}
