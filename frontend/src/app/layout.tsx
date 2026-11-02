import type { Metadata } from "next";

import { AppProviders } from "@/components/providers/app-providers";

import "./globals.css";

export const metadata: Metadata = {
  title: "Adfluence",
  description: "Ethiopia's creator marketplace: businesses find vetted local creators and work together with clear contracts in birr.",
  icons: {
    icon: [
      { url: "/brand/adfluence_favicon.ico" },
      { url: "/brand/adfluence_favicon.svg", type: "image/svg+xml" },
    ],
    apple: "/brand/adfluence_icon_128.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark h-full antialiased">
      <body className="min-h-full">
        <AppProviders>
          {children}
        </AppProviders>
      </body>
    </html>
  );
}
