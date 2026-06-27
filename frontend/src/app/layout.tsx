import type { Metadata } from "next";
import { cookies } from "next/headers";

import { AppProviders } from "@/components/providers/app-providers";

import "./globals.css";

export const metadata: Metadata = {
  title: "Adfluence",
  description: "Adfluence influencer marketing marketplace frontend",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const cookieStore = await cookies();
  const rawTheme = cookieStore.get("app_theme")?.value;
  const defaultTheme =
    rawTheme === "light" || rawTheme === "dark" || rawTheme === "system"
      ? rawTheme
      : "system";
  const themeScale = cookieStore.get("theme_scale")?.value ?? "100";

  return (
    <html lang="en" suppressHydrationWarning className="h-full antialiased">
      <body className="min-h-full">
        <AppProviders defaultTheme={defaultTheme} themeScale={themeScale}>
          {children}
        </AppProviders>
      </body>
    </html>
  );
}
