import type { Metadata } from "next";
import { Abyssinica_SIL, Bricolage_Grotesque, Inter, Noto_Sans_Ethiopic } from "next/font/google";

import { AppProviders } from "@/components/providers/app-providers";

import "./globals.css";

// Inter for reading, Bricolage Grotesque for headlines, Noto Sans Ethiopic so
// Amharic in chats and profiles renders properly, Abyssinica SIL for Amharic headlines.
const inter = Inter({ subsets: ["latin"], variable: "--font-inter", display: "swap" });
const bricolage = Bricolage_Grotesque({ subsets: ["latin"], variable: "--font-bricolage", display: "swap" });
const notoEthiopic = Noto_Sans_Ethiopic({ subsets: ["ethiopic"], variable: "--font-noto-ethiopic", display: "swap" });
const abyssinica = Abyssinica_SIL({ subsets: ["ethiopic"], weight: "400", variable: "--font-abyssinica", display: "swap", preload: false });

export const metadata: Metadata = {
  title: "Adfluence · Ethiopia's creator marketplace",
  description: "Ethiopia's creator marketplace: businesses find vetted local creators and work together with clear contracts in birr.",
  icons: {
    icon: [
      { url: "/brand/adfluence_favicon.ico", sizes: "any" },
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
    <html
      lang="en"
      className={`dark h-full antialiased ${inter.variable} ${bricolage.variable} ${notoEthiopic.variable} ${abyssinica.variable}`}
    >
      <body className="min-h-full">
        <AppProviders>
          {children}
        </AppProviders>
      </body>
    </html>
  );
}
